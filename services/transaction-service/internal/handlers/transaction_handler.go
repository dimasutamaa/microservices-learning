package handlers

import (
	"errors"
	"net/http"
	"transaction-service/internal/client"
	"transaction-service/internal/dto"
	"transaction-service/internal/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// POST /transactions/topup
func (h *TransactionHandler) Topup(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required header: X-User-ID"})
		return
	}

	var request dto.TopUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.transactionService.Topup(userID, request.Amount)
	if err != nil {
		switch {
		case errors.Is(err, client.ErrWalletServiceUnavailable):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Wallet service unavailable"})
		case errors.Is(err, client.ErrWalletService):
			c.JSON(http.StatusBadGateway, gin.H{"error": "Wallet service error"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error occured"})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
