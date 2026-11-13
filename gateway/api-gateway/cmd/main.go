package main

import (
	"api-gateway/internal/client"
	"api-gateway/internal/config"
	"api-gateway/internal/gateway"
	"api-gateway/internal/middleware"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found - using environment variables")
	}

	rdb, err := initRedis()
	if err != nil {
		log.Fatalf("[Redis] %v", err)
	}
	defer rdb.Close()

	router := gin.Default()
	router.Use(cors.Default())

	cfg := config.LoadGlobalConfig()

	authClient := client.NewAuthClient(cfg.AuthServiceBaseURL)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret, rdb, authClient)

	// Auth Service
	authProxy := gateway.NewProxy(cfg.AuthServiceBaseURL)
	router.POST("/api/v1/auth/register", authProxy)
	router.POST("/api/v1/auth/login", authProxy)
	router.POST("/api/v1/auth/logout", authMiddleware.LogoutAuth(), authProxy)

	router.Run(":8080")
}

func initRedis() (*redis.Client, error) {
	redisConfig := config.LoadRedisConfig()

	redisDB, err := strconv.Atoi(redisConfig.DB)
	if err != nil {
		return nil, fmt.Errorf("Invalid Redis DB index: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisDB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("Failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis")
	return rdb, nil
}
