package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:'user'"` // user, admin
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Document represents a document in the system
type Document struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	Title         string         `json:"title" gorm:"not null"`
	OriginalName  string         `json:"original_name" gorm:"not null"`
	FilePath      string         `json:"file_path" gorm:"not null"`
	FileSize      int64          `json:"file_size"`
	FileType      string         `json:"file_type"`
	Status        string         `json:"status" gorm:"default:'uploading'"` // uploading, processing, completed, failed
	ProcessedAt   *time.Time     `json:"processed_at"`
	VectorCount   int            `json:"vector_count" gorm:"default:0"`
	CollectionID  string         `json:"collection_id"`
	ErrorMessage  string         `json:"error_message"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// Conversation represents a RAG conversation
type Conversation struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Title     string         `json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Messages  []Message      `json:"messages" gorm:"foreignKey:ConversationID"`
}

// Message represents a single message in a conversation
type Message struct {
	ID             uint           `json:"id" gorm:"primarykey"`
	ConversationID uint           `json:"conversation_id" gorm:"not null"`
	Role           string         `json:"role" gorm:"not null"` // user, assistant
	Content        string         `json:"content" gorm:"type:text"`
	Sources        string         `json:"sources" gorm:"type:json"` // JSON array of source documents
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// DocumentChunk represents a chunk of a processed document
type DocumentChunk struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	DocumentID uint      `json:"document_id" gorm:"not null"`
	Document   Document  `json:"document" gorm:"foreignKey:DocumentID"`
	ChunkIndex int       `json:"chunk_index"`
	Content    string    `json:"content" gorm:"type:text"`
	VectorID   string    `json:"vector_id"` // ID in Milvus
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}