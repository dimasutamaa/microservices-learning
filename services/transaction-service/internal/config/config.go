package config

import (
	"fmt"
	"transaction-service/internal/utils"
)

type GlobalConfig struct {
	WalletSvcBaseUrl string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		WalletSvcBaseUrl: utils.GetEnvOrFatal("WALLET_SERVICE_BASE_URL"),
	}
}

func LoadDBConfig() *DBConfig {
	return &DBConfig{
		Host:     utils.GetEnv("DB_HOST", "localhost"),
		Port:     utils.GetEnv("DB_PORT", "5432"),
		User:     utils.GetEnv("DB_USER", "postgres"),
		Password: utils.GetEnv("DB_PASSWORD", "postgres"),
		DBName:   utils.GetEnv("DB_NAME", "mydb"),
		SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
	}
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host,
		c.User,
		c.Password,
		c.DBName,
		c.Port,
		c.SSLMode,
	)
}
