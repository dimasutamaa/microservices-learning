package utils

import (
	"log"
	"os"
)

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Environment variable for %s not found, use default value", key)
	return defaultVal
}

func GetEnvOrFatal(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but missing", key)
	}
	return value
}
