package handlers

import (
	"auth-service/internal/dto"
	"auth-service/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewUserHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// POST /register
func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService, err := h.authService.Register(&request)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occured"})
		return
	}

	c.JSON(http.StatusCreated, authService)
}

// POST /login
func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(&request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occured"})
		return
	}

	c.JSON(http.StatusOK, token)
}

// POST /logout
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenID := c.GetHeader("X-Token-ID")
	userID := c.GetHeader("X-User-ID")

	if err := h.authService.Logout(tokenID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token is already revoked or invalid"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not log out"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
