package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"ragwitheitno/config"
	"ragwitheitno/internal/auth"
	"ragwitheitno/internal/database"
	"ragwitheitno/internal/handlers"
	"ragwitheitno/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connections
	db, err := database.InitMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	redisClient, err := database.InitRedis(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	milvusClient, err := database.InitMilvus(cfg.MilvusAddr)
	if err != nil {
		log.Fatal("Failed to connect to Milvus:", err)
	}

	// Initialize services
	authService := auth.NewService(redisClient, cfg.JWTSecret)
	
	// Initialize handlers
	userHandler := handlers.NewUserHandler(db, authService)
	documentHandler := handlers.NewDocumentHandler(db, redisClient, milvusClient)
	ragHandler := handlers.NewRAGHandler(db, milvusClient)

	// Setup Gin router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/register", userHandler.Register)
			auth.POST("/logout", middleware.AuthRequired(authService), userHandler.Logout)
			auth.GET("/profile", middleware.AuthRequired(authService), userHandler.GetProfile)
		}

		// User management routes (admin only)
		users := api.Group("/users")
		users.Use(middleware.AuthRequired(authService))
		users.Use(middleware.AdminRequired())
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Document management routes
		documents := api.Group("/documents")
		documents.Use(middleware.AuthRequired(authService))
		{
			documents.GET("", documentHandler.ListDocuments)
			documents.POST("", documentHandler.UploadDocument)
			documents.GET("/:id", documentHandler.GetDocument)
			documents.DELETE("/:id", documentHandler.DeleteDocument)
			documents.GET("/:id/status", documentHandler.GetProcessingStatus)
		}

		// RAG query routes
		rag := api.Group("/rag")
		rag.Use(middleware.AuthRequired(authService))
		{
			rag.POST("/query", ragHandler.Query)
			rag.GET("/conversations", ragHandler.GetConversations)
			rag.GET("/conversations/:id", ragHandler.GetConversation)
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}