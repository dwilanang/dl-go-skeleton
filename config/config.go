package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName       string
	AppPort       string
	DBDriver      string
	DBHost        string
	DBName        string
	DBPassword    string
	DBUser        string
	DBPort        string
	JWTSecret     string
	JWTExpiration string
	JWTType       string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		AppName:       getEnv("APP_NAME", "payslip-service"),
		AppPort:       getEnv("APP_PORT", "8000"),
		DBDriver:      getEnv("DB_DRIVER", "postgres"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", ""),
		DBName:        getEnv("DB_NAME", ""),
		DBUser:        getEnv("DB_USER", ""),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiration: getEnv("JWT_EXPIRATION", "0"),
		JWTType:       getEnv("JWT_TYPE", "bearer"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
