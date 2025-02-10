package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPassword             string
	DBName                 string
	SMTPHost               string
	SMTPPort               string
	SMTPUsername           string
	SMTPPassword           string
	SMTPEmail              string
	EncryptKey             string
	EncryptIv              string
	FrontendUrl            string
	RedisHost              string
	RedisPort              string
	RedisPassword          string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		DBHost:                 getENV("DB_HOST", "localhost"),
		DBPort:                 getENV("DB_PORT", "5432"),
		DBUser:                 getENV("DB_USER", "postgres"),
		DBPassword:             getENV("DB_PASSWORD", "postgres"),
		DBName:                 getENV("DB_NAME", "postgres"),
		SMTPHost:               getENV("SMTP_HOST", ""),
		SMTPPort:               getENV("SMTP_PORT", ""),
		SMTPUsername:           getENV("SMTP_USERNAME", ""),
		SMTPPassword:           getENV("SMTP_PASSWORD", ""),
		SMTPEmail:              getENV("SMTP_EMAIL", ""),
		EncryptKey:             getENV("ENCRYPT_KEY", ""),
		EncryptIv:              getENV("ENCRYPT_IV", ""),
		FrontendUrl:            getENV("FRONTEND_URL", "http://localhost:3000"),
		RedisHost:              getENV("REDIS_HOST", "127.0.0.1"),
		RedisPort:              getENV("REDIS_PORT", "6379"),
		RedisPassword:          getENV("REDIS_PASSWORD", ""),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
		JWTSecret:              getENV("JWT_SECRET", "secret"),
	}
}

func getENV(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
