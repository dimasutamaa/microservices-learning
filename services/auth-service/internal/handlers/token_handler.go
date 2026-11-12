package handlers

import (
	"auth-service/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TokenHandler struct {
	tokenService *services.AccessTokenService
}

func NewTokenHandler(tokenService *services.AccessTokenService) *TokenHandler {
	return &TokenHandler{tokenService: tokenService}
}

// GET /token/{tokenID}/usability
func (h *TokenHandler) CheckTokenUsability(c *gin.Context) {
	tokenID := c.Param("tokenID")

	valid, err := h.tokenService.IsTokenUsable(tokenID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, "Invalid token or expired")
			return
		}

		c.JSON(http.StatusInternalServerError, "Unexpected error occured")
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": valid})
}
