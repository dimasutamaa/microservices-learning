package handlers

import (
	"errors"
	"net/http"
	"wallet-service/internal/client"
	"wallet-service/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

// POST /wallet/create
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	userID := c.Request.Header.Get("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required header: X-User-ID"})
		return
	}

	wallet, err := h.walletService.CreateWallet(userID)
	if err != nil {
		switch {
		case errors.Is(err, client.ErrAuthServiceUnavailable):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return

		case errors.Is(err, client.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		case errors.Is(err, client.ErrAuthService):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		case errors.Is(err, gorm.ErrDuplicatedKey):
			c.JSON(http.StatusConflict, gin.H{"error": "User with this userID already has a wallet"})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occured"})
			return
		}
	}

	c.JSON(http.StatusCreated, wallet)
}

// GET /wallet/balance
func (h *WalletHandler) GetWalletBalance(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required header: X-User-ID"})
		return
	}

	wallet, err := h.walletService.GetWalletBalance(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User wallet not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occured"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}
