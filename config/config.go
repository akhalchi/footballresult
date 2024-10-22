package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Configuration stores all environment variables
type Configuration struct {
	FootballDataToken    string
	FootballDataURL      string
	TelegramBotToken     string
	TelegramChannelID    string
	TelegramLogChannelID string
	DBUser               string
	DBPass               string
	DBName               string
	DBHost               string
	DBPort               string
}

// Global variable to hold the configuration
var Load Configuration

// LoadConfig loads environment variables from the .env file
func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	} else {
		log.Printf("env file is loaded")
	}

	Load = Configuration{
		FootballDataToken:    getEnv("FOOTBALL_DATA_TOKEN"),
		FootballDataURL:      getEnv("FOOTBALL_DATA_URL"),
		TelegramBotToken:     getEnv("TELEGRAM_BOT_TOKEN"),
		TelegramChannelID:    getEnv("TELEGRAM_CHANNEL_ID"),
		TelegramLogChannelID: getEnv("TELEGRAM_LOG_CHANNEL_ID"),
		DBUser:               getEnv("DB_USER"),
		DBPass:               getEnv("DB_PASS"),
		DBName:               getEnv("DB_NAME"),
		DBHost:               getEnv("DB_HOST"),
		DBPort:               getEnv("DB_PORT"),
	}
}

// getEnv fetches the environment variable by key
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s not set in .env file", key)
	} else {
		log.Print("variable is loaded: ", key)
	}
	return value
}
