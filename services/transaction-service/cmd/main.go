package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/transactions/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "transaction-service",
		})
	})

	router.Run(":8083")
}
