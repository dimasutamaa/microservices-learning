package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func createProxy(target string) gin.HandlerFunc {
	targetURL, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "api-gateway",
		})
	})

	router.Any("/auth/*path", createProxy("http://localhost:8081"))
	router.Any("/wallets/*path", createProxy("http://localhost:8082"))
	router.Any("/transactions/*path", createProxy("http://localhost:8083"))

	router.Run(":8080")
}
