package routes

import (
	"auth-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, authHandler *handlers.AuthHandler, tokenHandler *handlers.TokenHandler) {
	api := router.Group("/auth")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/logout", authHandler.Logout)

		api.GET("/token/:tokenID/usability", tokenHandler.CheckTokenUsability)
	}
}
