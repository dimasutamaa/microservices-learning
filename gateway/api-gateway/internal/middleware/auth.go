package middleware

import (
	"api-gateway/internal/client"
	"api-gateway/internal/utils"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	jwtSecret  []byte
	rdb        *redis.Client
	authClient *client.AuthClient
}

func NewAuthMiddleware(jwtSecret []byte, rdb *redis.Client, authClient *client.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:  jwtSecret,
		rdb:        rdb,
		authClient: authClient,
	}
}

func (m *AuthMiddleware) authenticate(c *gin.Context) (*utils.Claims, bool) {
	parts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		unauthorized(c, "Missing authorization header")
		return nil, false
	}

	claims, err := utils.ParseToken(parts[1], m.jwtSecret)
	if err != nil {
		unauthorized(c, "Invalid or expired token")
		return nil, false
	}

	ctx := c.Request.Context()

	// Redis cache hit
	if m.rdb != nil {
		cachedUserID, err := m.rdb.Get(ctx, cacheKey(claims.TokenID)).Result()
		if err == nil {
			if cachedUserID != claims.UserID {
				unauthorized(c, "Invalid token")
				return nil, false
			}
			return claims, true
		}
	}

	// Redis cache miss -> Auth Service
	authCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	valid, err := m.authClient.IsTokenUsable(authCtx, claims.TokenID)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Authorization service unavailable",
		})
		c.Abort()
		return nil, false
	}

	if !valid {
		unauthorized(c, "Access token has been revoked or expired")
		return nil, false
	}

	if m.rdb != nil {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			_ = m.rdb.Set(ctx, cacheKey(claims.TokenID), claims.UserID, ttl).Err()
		}
	}

	return claims, true
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := m.authenticate(c)
		if !ok {
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) LogoutAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := m.authenticate(c)
		if !ok {
			return
		}

		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-Token-ID", claims.TokenID)
		c.Next()
	}
}

func cacheKey(tokenID string) string {
	return "access-token:" + tokenID
}

func unauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": msg,
	})
}
