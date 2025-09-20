package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	MySQLDSN      string
	RedisAddr     string
	RedisPassword string
	MilvusAddr    string
	JWTSecret     string
	JWTExpiry     time.Duration
	EinoConfig    EinoConfig
}

type EinoConfig struct {
	ModelName    string
	EmbeddingDim int
	MaxTokens    int
}

func Load() *Config {
	jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	embeddingDim, _ := strconv.Atoi(getEnv("EMBEDDING_DIM", "1536"))
	maxTokens, _ := strconv.Atoi(getEnv("MAX_TOKENS", "4096"))

	return &Config{
		Port:          getEnv("PORT", "8080"),
		MySQLDSN:      getEnv("MYSQL_DSN", "root:password@tcp(localhost:3306)/ragdb?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		MilvusAddr:    getEnv("MILVUS_ADDR", "localhost:19530"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiry:     time.Duration(jwtExpiryHours) * time.Hour,
		EinoConfig: EinoConfig{
			ModelName:    getEnv("MODEL_NAME", "gpt-3.5-turbo"),
			EmbeddingDim: embeddingDim,
			MaxTokens:    maxTokens,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}