package main

import (
	"api-gateway/internal/gateway"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	// Auth Service
	authProxy := gateway.NewProxy("http://localhost:8081")
	router.Any("/api/v1/auth/*path", authProxy)

	router.Run(":8080")
}
