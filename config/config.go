package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DBPath             string
	DatabaseURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	JWTSecret          string
	JWTExpiration      int
	ServerPort         string
	FrontendURL        string
}

var AppConfig Config

func InitConfig() {
	// Load environment variables from .env (if present)
	_ = godotenv.Load()

	AppConfig = Config{
		DBPath:             getEnv("DB_PATH", "gofiber.db"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", "secret"),
		JWTExpiration:      getEnvAsInt("JWT_EXPIRATION", 72),
		ServerPort:         getEnv("SERVER_PORT", "8000"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
