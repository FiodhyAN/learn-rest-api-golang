package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPEmail    string
	EncryptKey   string
	EncryptIv    string
	FrontendUrl  string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		DBHost:       getENV("DB_HOST", "localhost"),
		DBPort:       getENV("DB_PORT", "5432"),
		DBUser:       getENV("DB_USER", "postgres"),
		DBPassword:   getENV("DB_PASSWORD", "postgres"),
		DBName:       getENV("DB_NAME", "postgres"),
		SMTPHost:     getENV("SMTP_HOST", ""),
		SMTPPort:     getENV("SMTP_PORT", ""),
		SMTPUsername: getENV("SMTP_USERNAME", ""),
		SMTPPassword: getENV("SMTP_PASSWORD", ""),
		SMTPEmail:    getENV("SMTP_EMAIL", ""),
		EncryptKey:   GenerateRandomString(32),
		EncryptIv:    GenerateRandomString(16),
		FrontendUrl:  getENV("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getENV(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
