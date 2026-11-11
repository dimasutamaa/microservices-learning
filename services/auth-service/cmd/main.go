package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/routes"
	"auth-service/internal/services"
	"auth-service/utils"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found - using environment variables")
	}

	cfg := config.LoadGlobalConfig()

	db, err := initDatabase()
	if err != nil {
		log.Fatalf("[Database] %v", err)
	}

	rdb, err := initRedis()
	if err != nil {
		log.Fatalf("[Redis] %v", err)
	}
	defer rdb.Close()

	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewAccessTokenRepository(db)
	userSvc := services.NewUserService(userRepo, tokenRepo, rdb, cfg)
	userHandler := handlers.NewUserHandler(userSvc)

	router := gin.Default()
	routes.SetupRouter(router, userHandler)

	addr := utils.GetEnv("PORT", "8081")
	if err := router.Run(":" + addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initDatabase() (*gorm.DB, error) {
	dbConfig := config.LoadDBConfig()

	db, err := gorm.Open(postgres.Open(dbConfig.DSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to PostgreSQL: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.AccessToken{}); err != nil {
		return nil, fmt.Errorf("Failed to migrate database: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
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
