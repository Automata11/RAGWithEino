package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"gorm.io/gorm"
	"ragwitheitno/internal/middleware"
	"ragwitheitno/internal/models"
)

type DocumentHandler struct {
	db           *gorm.DB
	redisClient  *redis.Client
	milvusClient client.Client
}

func NewDocumentHandler(db *gorm.DB, redisClient *redis.Client, milvusClient client.Client) *DocumentHandler {
	return &DocumentHandler{
		db:           db,
		redisClient:  redisClient,
		milvusClient: milvusClient,
	}
}

// ListDocuments returns user's documents
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	var documents []models.Document
	query := h.db.Where("user_id = ?", userID)

	// Check if user is admin, show all documents
	role, _ := middleware.GetCurrentUserRole(c)
	if role == "admin" {
		query = h.db.Preload("User")
	}

	if err := query.Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// UploadDocument handles document upload
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := time.Now().Format("20060102-150405") + "_" + header.Filename
	filePath := filepath.Join(uploadDir, filename)

	// Save file to disk
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Get file info
	fileInfo, err := dst.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	// Determine file type
	fileType := getFileType(header.Filename)
	if fileType == "unsupported" {
		// Clean up uploaded file
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
		return
	}

	// Create document record
	title := c.PostForm("title")
	if title == "" {
		title = header.Filename
	}

	document := models.Document{
		UserID:       userID,
		Title:        title,
		OriginalName: header.Filename,
		FilePath:     filePath,
		FileSize:     fileInfo.Size(),
		FileType:     fileType,
		Status:       "uploading",
	}

	if err := h.db.Create(&document).Error; err != nil {
		// Clean up uploaded file
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document record"})
		return
	}

	// Start background processing
	go h.processDocument(document.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Document uploaded successfully",
		"document": document,
	})
}

// GetDocument returns a specific document
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var document models.Document
	query := h.db.Where("id = ? AND user_id = ?", id, userID)

	// Check if user is admin
	role, _ := middleware.GetCurrentUserRole(c)
	if role == "admin" {
		query = h.db.Where("id = ?", id).Preload("User")
	}

	if err := query.First(&document).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"document": document})
}

// DeleteDocument deletes a document
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var document models.Document
	query := h.db.Where("id = ? AND user_id = ?", id, userID)

	// Check if user is admin
	role, _ := middleware.GetCurrentUserRole(c)
	if role == "admin" {
		query = h.db.Where("id = ?", id)
	}

	if err := query.First(&document).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Delete file from disk
	if err := os.Remove(document.FilePath); err != nil {
		// Log error but continue with database deletion
	}

	// Delete from database (this will also delete chunks due to foreign key constraints)
	if err := h.db.Delete(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// GetProcessingStatus returns document processing status
func (h *DocumentHandler) GetProcessingStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Get status from database first
	var document models.Document
	if err := h.db.First(&document, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Check Redis for more detailed status if processing
	var detailedStatus string
	if document.Status == "processing" {
		ctx := c.Request.Context()
		key := "doc_status:" + strconv.FormatUint(id, 10)
		status, err := h.redisClient.Get(ctx, key).Result()
		if err == nil {
			detailedStatus = status
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          document.Status,
		"detailed_status": detailedStatus,
		"processed_at":    document.ProcessedAt,
		"vector_count":    document.VectorCount,
		"error_message":   document.ErrorMessage,
	})
}

// processDocument processes a document in the background
func (h *DocumentHandler) processDocument(documentID uint) {
	// Update status to processing
	h.db.Model(&models.Document{}).Where("id = ?", documentID).Update("status", "processing")

	// Set processing status in Redis
	ctx := h.redisClient.Context()
	statusKey := "doc_status:" + strconv.FormatUint(uint64(documentID), 10)
	h.redisClient.Set(ctx, statusKey, "extracting_text", time.Hour)

	// TODO: Implement actual document processing
	// This would include:
	// 1. Text extraction based on file type
	// 2. Text chunking
	// 3. Embedding generation using Eino
	// 4. Vector storage in Milvus
	// 5. Metadata storage in database

	// For now, simulate processing
	time.Sleep(5 * time.Second) // Simulate processing time

	// Update status to completed
	now := time.Now()
	h.db.Model(&models.Document{}).Where("id = ?", documentID).Updates(map[string]interface{}{
		"status":       "completed",
		"processed_at": &now,
		"vector_count": 10, // Simulated
	})

	// Remove processing status from Redis
	h.redisClient.Del(ctx, statusKey)
}

// getFileType determines file type based on extension
func getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "pdf"
	case ".txt":
		return "text"
	case ".doc", ".docx":
		return "document"
	case ".md":
		return "markdown"
	default:
		return "unsupported"
	}
}