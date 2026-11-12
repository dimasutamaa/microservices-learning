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
