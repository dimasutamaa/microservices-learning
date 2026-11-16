package handlers

import (
	"errors"
	"net/http"
	"wallet-service/internal/dto"
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
	var request dto.WalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.CreateWallet(request.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this userID already has a wallet"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occured"})
		return
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
