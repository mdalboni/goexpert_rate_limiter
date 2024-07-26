package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	IPLimit        int
	IPDuration     time.Duration
	APIKeyLimit    int
	APIKeyDuration time.Duration
	RedisAddress   string
	RedisPassword  string
}

func getEnvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvStr(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file, will use environment variables")
	}
	ipLimit := getEnvInt("IP_LIMIT", 10)
	ipDuration := getEnvInt("IP_LIMIT_DURATION", 1)
	apiKeyLimit := getEnvInt("API_KEY_LIMIT", 100)
	apiKeyDuration := getEnvInt("API_KEY_LIMIT_DURATION", 1)
	redisAddress := getEnvStr("REDIS_ADDRESS", "localhost:6379")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	return Config{
		IPLimit:        ipLimit,
		IPDuration:     time.Duration(ipDuration) * time.Second,
		APIKeyLimit:    apiKeyLimit,
		APIKeyDuration: time.Duration(apiKeyDuration) * time.Second,
		RedisAddress:   redisAddress,
		RedisPassword:  redisPassword,
	}
}
