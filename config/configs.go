package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found,using system environment variables")
	}

	return &Config{
		ServerPort: getenv("SERVER_PORT", "8888"),
		DBHost:     getenv("DB_HOST", "db"),
		DBPort:     getenv("DB_PORT", "3306"),
		DBUser:     getenv("DB_USER", "appuser"),
		DBPassword: getenv("DB_PASSWORD", "123456"),
		DBName:     getenv("DB_NAME", "steam"),
		JWTSecret:  getenv("JWT_SECRET", "key"),
	}
}

func getenv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
