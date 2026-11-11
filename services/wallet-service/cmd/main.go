package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/wallets/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"services": "wallet-service",
		})
	})

	router.Run(":8082")
}
