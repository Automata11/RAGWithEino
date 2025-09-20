package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	redisClient *redis.Client
	jwtSecret   []byte
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func NewService(redisClient *redis.Client, jwtSecret string) *Service {
	return &Service{
		redisClient: redisClient,
		jwtSecret:   []byte(jwtSecret),
	}
}

// HashPassword hashes a password using bcrypt
func (s *Service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword verifies a password against its hash
func (s *Service) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken creates a new JWT token
func (s *Service) GenerateToken(userID uint, username, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	// Store token in Redis with expiration
	ctx := context.Background()
	key := fmt.Sprintf("token:%d", userID)
	err = s.redisClient.Set(ctx, key, tokenString, 24*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store token in Redis: %v", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check if token exists in Redis
	ctx := context.Background()
	key := fmt.Sprintf("token:%d", claims.UserID)
	storedToken, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("token not found or expired")
	} else if err != nil {
		return nil, fmt.Errorf("Redis error: %v", err)
	}

	if storedToken != tokenString {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// RevokeToken removes a token from Redis
func (s *Service) RevokeToken(userID uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("token:%d", userID)
	return s.redisClient.Del(ctx, key).Err()
}

// SetProcessingStatus sets document processing status in Redis
func (s *Service) SetProcessingStatus(documentID uint, status string) error {
	ctx := context.Background()
	key := fmt.Sprintf("doc_status:%d", documentID)
	return s.redisClient.Set(ctx, key, status, 1*time.Hour).Err()
}

// GetProcessingStatus gets document processing status from Redis
func (s *Service) GetProcessingStatus(documentID uint) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("doc_status:%d", documentID)
	status, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "unknown", nil
	}
	return status, err
}