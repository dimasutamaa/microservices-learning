package config

import (
	"api-gateway/internal/utils"
)

type GlobalConfig struct {
	AuthServiceBaseURL string
	JWTSecret          []byte
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       string
}

func LoadGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		AuthServiceBaseURL: utils.GetEnvOrFatal("AUTH_SERVICE_BASE_URL"),
		JWTSecret:          []byte(utils.GetEnvOrFatal("JWT_SECRET")),
	}
}

func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     utils.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       utils.GetEnv("REDIS_DB", "0"),
	}
}
