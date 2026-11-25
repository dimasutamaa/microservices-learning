package main

import (
	"fmt"
	"log"
	"wallet-service/internal/config"
	"wallet-service/internal/handlers"
	"wallet-service/internal/messaging"
	"wallet-service/internal/models"
	"wallet-service/internal/repositories"
	"wallet-service/internal/routes"
	"wallet-service/internal/services"
	"wallet-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found - using environment variables")
	}

	db, err := initDatabase()
	if err != nil {
		log.Fatalf("[Database] %v", err)
	}

	walletRepo := repositories.NewWalletRepository(db)
	walletSvc := services.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletSvc)

	rabbitMQ, err := messaging.NewRabbitMQ(utils.GetEnvOrFatal("RABBITMQ_URL"), walletSvc)
	if err != nil {
		log.Fatalf("[RabbitMQ] %v", err)
	}
	defer rabbitMQ.Close()

	router := gin.Default()
	routes.SetupRouter(router, walletHandler)

	// Process transfer requests in the background
	go func() {
		if err := rabbitMQ.ConsumeTransferCommands(); err != nil {
			log.Fatalf("[RabbitMQ] transfer command consumer stopped: %v", err)
		}
	}()

	addr := utils.GetEnv("PORT", "8082")
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

	if err := db.AutoMigrate(&models.Wallet{}); err != nil {
		return nil, fmt.Errorf("Failed to migrate database: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}
