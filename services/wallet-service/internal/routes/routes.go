package routes

import (
	"wallet-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, walletHandler *handlers.WalletHandler) {
	api := router.Group("/wallet")
	{
		api.POST("/create", walletHandler.CreateWallet)
		api.GET("/balance", walletHandler.GetWalletBalance)
	}
}
