package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"gorm.io/gorm"
	"ragwitheitno/internal/middleware"
	"ragwitheitno/internal/models"
)

type RAGHandler struct {
	db           *gorm.DB
	milvusClient client.Client
}

type QueryRequest struct {
	Question       string `json:"question" binding:"required"`
	ConversationID *uint  `json:"conversation_id,omitempty"`
}

type QueryResponse struct {
	Answer         string                 `json:"answer"`
	Sources        []DocumentSource       `json:"sources"`
	ConversationID uint                   `json:"conversation_id"`
	MessageID      uint                   `json:"message_id"`
}

type DocumentSource struct {
	DocumentID   uint   `json:"document_id"`
	DocumentName string `json:"document_name"`
	ChunkID      uint   `json:"chunk_id"`
	Content      string `json:"content"`
	Similarity   float32 `json:"similarity"`
}

func NewRAGHandler(db *gorm.DB, milvusClient client.Client) *RAGHandler {
	return &RAGHandler{
		db:           db,
		milvusClient: milvusClient,
	}
}

// Query handles RAG query requests
func (h *RAGHandler) Query(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get or create conversation
	var conversation models.Conversation
	if req.ConversationID != nil {
		if err := h.db.Where("id = ? AND user_id = ?", *req.ConversationID, userID).First(&conversation).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
			return
		}
	} else {
		// Create new conversation
		conversation = models.Conversation{
			UserID: userID,
			Title:  truncateString(req.Question, 50),
		}
		if err := h.db.Create(&conversation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
			return
		}
	}

	// Save user message
	userMessage := models.Message{
		ConversationID: conversation.ID,
		Role:           "user",
		Content:        req.Question,
	}
	if err := h.db.Create(&userMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// TODO: Implement actual RAG query using Eino
	// This would include:
	// 1. Generate embedding for the question using Eino
	// 2. Search similar vectors in Milvus
	// 3. Retrieve relevant document chunks
	// 4. Generate answer using Eino with retrieved context
	// 5. Return structured response

	// For now, simulate RAG response
	answer, sources := h.simulateRAGQuery(req.Question, userID)

	// Save assistant message
	assistantMessage := models.Message{
		ConversationID: conversation.ID,
		Role:           "assistant",
		Content:        answer,
		Sources:        marshalSources(sources), // Convert to JSON string
	}
	if err := h.db.Create(&assistantMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save response"})
		return
	}

	response := QueryResponse{
		Answer:         answer,
		Sources:        sources,
		ConversationID: conversation.ID,
		MessageID:      assistantMessage.ID,
	}

	c.JSON(http.StatusOK, response)
}

// GetConversations returns user's conversations
func (h *RAGHandler) GetConversations(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	var conversations []models.Conversation
	if err := h.db.Where("user_id = ?", userID).Order("updated_at DESC").Find(&conversations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

// GetConversation returns a specific conversation with messages
func (h *RAGHandler) GetConversation(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	var conversation models.Conversation
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Messages").First(&conversation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversation": conversation})
}

// simulateRAGQuery simulates a RAG query response
func (h *RAGHandler) simulateRAGQuery(question string, userID uint) (string, []DocumentSource) {
	// Simulate retrieving relevant documents
	var documents []models.Document
	h.db.Where("user_id = ? AND status = 'completed'", userID).Limit(3).Find(&documents)

	var sources []DocumentSource
	for i, doc := range documents {
		sources = append(sources, DocumentSource{
			DocumentID:   doc.ID,
			DocumentName: doc.Title,
			ChunkID:      uint(i + 1),
			Content:      "This is a simulated chunk from " + doc.Title + " that might be relevant to the question.",
			Similarity:   0.85 - float32(i)*0.1,
		})
	}

	// Simulate answer generation
	answer := "Based on the available documents, I can provide the following response to your question: \"" + question + "\". "
	if len(sources) > 0 {
		answer += "This information is derived from the uploaded documents including " + documents[0].Title
		if len(documents) > 1 {
			answer += " and others"
		}
		answer += ". Please note that this is a simulated response for demonstration purposes."
	} else {
		answer += "However, I don't have any relevant documents to provide a specific answer. Please upload some documents first."
	}

	return answer, sources
}

// truncateString truncates a string to the specified length
func truncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length] + "..."
}

// marshalSources converts sources to JSON string
func marshalSources(sources []DocumentSource) string {
	// In a real implementation, you would use json.Marshal
	// For now, return a simple representation
	if len(sources) == 0 {
		return "[]"
	}
	return "[{\"simulated\": \"sources\"}]" // Placeholder
}