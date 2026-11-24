package routes

import (
	"transaction-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, transactionHandler *handlers.TransactionHandler) {
	api := router.Group("/transactions")
	{
		api.POST("/topup", transactionHandler.Topup)
	}
}
