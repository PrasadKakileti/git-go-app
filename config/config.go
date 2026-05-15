package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string

	// Email (SMTP or SendGrid)
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPass       string
	SendGridAPIKey string

	// Job APIs
	// Get a free key at: https://rapidapi.com/letscrape-6bRBa3QguO5/api/jsearch
	// Free tier: 200 requests/month — covers LinkedIn, Indeed, Glassdoor, ZipRecruiter
	JSearchAPIKey string

	// App
	ServerPort string
	AppBaseURL string // e.g. http://localhost:8080 — used for email links

	// Admin credentials (loaded from env, NOT hardcoded)
	AdminEmail    string
	AdminPassword string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "job_portal"),

		SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPass:       getEnv("SMTP_PASSWORD", ""),
		SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),

		JSearchAPIKey: getEnv("JSEARCH_API_KEY", ""),

		ServerPort: getEnv("SERVER_PORT", "8080"),
		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),

		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@jobhub.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
