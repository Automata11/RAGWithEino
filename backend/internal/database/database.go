package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ragwitheitno/internal/models"
)

// InitMySQL initializes MySQL database connection
func InitMySQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Auto-migrate the models
	err = db.AutoMigrate(
		&models.User{},
		&models.Document{},
		&models.Conversation{},
		&models.Message{},
		&models.DocumentChunk{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	// Create default admin user if not exists
	var adminUser models.User
	if db.Where("username = ?", "admin").First(&adminUser).Error == gorm.ErrRecordNotFound {
		// Hash password (you should implement proper password hashing)
		adminUser = models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:     "admin",
			IsActive: true,
		}
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
		}
	}

	return db, nil
}

// InitRedis initializes Redis client
func InitRedis(addr, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}

// InitMilvus initializes Milvus client
func InitMilvus(addr string) (client.Client, error) {
	c, err := client.NewGrpcClient(context.Background(), addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Milvus: %v", err)
	}

	// Note: Collection creation simplified for now
	// In a real implementation, you would create the proper schema here
	log.Printf("Connected to Milvus at %s", addr)

	return c, nil
}