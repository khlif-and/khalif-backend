package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	AppEnv     string

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string

	JWTSecret    string
	JWTExpHours  int
	RefreshTokenSecret string
	RefreshTokenExpDays int

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Brevo Email Service
	BrevoAPIKey       string
	BrevoAPIURL       string
	BrevoSenderEmail  string
	BrevoSenderName   string
	EmailTemplatePath string

	// Meilisearch
	MeilisearchHost   string
	MeilisearchAPIKey string
}

func LoadConfig() *Config {
	// Load .env file if available, otherwise ignore error and return system env vars
	_ = godotenv.Load()

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "khalif_backend"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:   getEnv("JWT_SECRET", "default_secret"),
		JWTExpHours: getEnvAsInt("JWT_EXP_HOURS", 24),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "default_refresh_secret"),
		RefreshTokenExpDays: getEnvAsInt("REFRESH_TOKEN_EXP_DAYS", 7),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		BrevoAPIKey:       getEnv("BREVO_API_KEY", ""),
		BrevoAPIURL:       getEnv("BREVO_API_URL", "https://api.brevo.com/v3/smtp/email"),
		BrevoSenderEmail:  getEnv("BREVO_SENDER_EMAIL", "noreply@khalifapp.com"),
		BrevoSenderName:   getEnv("BREVO_SENDER_NAME", "Khalif App"),
		EmailTemplatePath: getEnv("EMAIL_TEMPLATE_PATH", "templates/email"),

		MeilisearchHost:   getEnv("MEILISEARCH_HOST", "http://localhost:7700"),
		MeilisearchAPIKey: getEnv("MEILISEARCH_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}
