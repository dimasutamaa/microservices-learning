package utils

import (
	"auth-service/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	TokenID string `json:"id"`
	UserID  string `json:"subject"`
	jwt.RegisteredClaims
}

func GenerateToken(at *models.AccessToken, secret []byte, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		TokenID: at.ID,
		UserID:  at.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "microservices-learning",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
