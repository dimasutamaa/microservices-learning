package main

import (
	"fmt"
	"log"
	"transaction-service/internal/client"
	"transaction-service/internal/config"
	"transaction-service/internal/handlers"
	"transaction-service/internal/messaging"
	"transaction-service/internal/models"
	"transaction-service/internal/repositories"
	"transaction-service/internal/routes"
	"transaction-service/internal/services"
	"transaction-service/internal/utils"

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

	cfg := config.LoadGlobalConfig()
	trxRepo := repositories.NewTransactionRepository(db)

	rabbitMQ, err := messaging.NewRabbitMQ(cfg.RabbitMQURL, trxRepo)
	if err != nil {
		log.Fatalf("[RabbitMQ] %v", err)
	}
	defer rabbitMQ.Close()

	walletClient := client.NewWalletClient(cfg.WalletSvcBaseUrl)
	trxService := services.NewTransactionService(trxRepo, walletClient, rabbitMQ)
	trxHandler := handlers.NewTransactionHandler(trxService)

	router := gin.Default()
	routes.SetupRouter(router, trxHandler)

	// Listen for wallet results in the background
	go func() {
		if err := rabbitMQ.ConsumeTransferResults(); err != nil {
			log.Fatalf("[RabbitMQ] transfer result consumer stopped: %v", err)
		}
	}()

	addr := utils.GetEnv("PORT", "8083")
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

	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		return nil, fmt.Errorf("Failed to migrate database: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}
