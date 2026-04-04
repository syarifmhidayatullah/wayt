package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  string
	AppEnv   string
	DB       DBConfig
	Internal InternalConfig
	QR       QRConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type InternalConfig struct {
	APIKey string
}

type QRConfig struct {
	StoragePath   string
	BaseURL       string
	PublicBaseURL string
	ExpiredHours  int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	expiredHours, err := strconv.Atoi(getEnv("QR_EXPIRED_HOURS", "24"))
	if err != nil {
		expiredHours = 24
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "wayt"),
		},
		Internal: InternalConfig{
			APIKey: getEnv("INTERNAL_API_KEY", ""),
		},
		QR: QRConfig{
			StoragePath:   getEnv("QR_STORAGE_PATH", "./storage/qr"),
			BaseURL:       getEnv("QR_BASE_URL", "http://localhost:8080/storage/qr"),
			PublicBaseURL: getEnv("PUBLIC_BASE_URL", "http://localhost:8080"),
			ExpiredHours:  expiredHours,
		},
	}, nil
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
