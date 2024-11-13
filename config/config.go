package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ENV      string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBEngine string
	Port     int
	DBName   string

	AllowedOrigins []string
	AllowedHeaders []string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	allowed := os.Getenv("CORS_ALLOWED_ORIGINS")

	return &Config{
		ENV:            os.Getenv("ENV"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBName:         os.Getenv("DB_NAME"),
		DBUser:         os.Getenv("DB_USER"),
		DBPass:         os.Getenv("DB_PASS"),
		DBEngine:       os.Getenv("DB_ENGINE"),
		Port:           port,
		AllowedOrigins: strings.Split(allowed, ","),
		AllowedHeaders: strings.Split(os.Getenv("CORS_ALLOWED_HEADERS"), ","),
	}
}
