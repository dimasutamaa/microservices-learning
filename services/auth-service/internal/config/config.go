package config

import (
	"auth-service/internal/utils"
	"fmt"
	"log"
	"time"
)

type GlobalConfig struct {
	TokenTTL  time.Duration
	JWTSecret []byte
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       string
}

func LoadGlobalConfig() *GlobalConfig {
	secret := utils.GetEnv("JWT_SECRET", "")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return &GlobalConfig{
		JWTSecret: []byte(secret),
		TokenTTL:  15 * time.Minute,
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

func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     utils.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       utils.GetEnv("REDIS_DB", "0"),
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
