package config

import (
	"os"
	"strconv"
	"time"

	"github.com/utkarsh/unfollow-tracker/internal/database"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Scraper  ScraperConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type ScraperConfig struct {
	DelayMs           int
	MaxConcurrent     int
	ProxyPoolURL      string
	NotificationDelay time.Duration
	InstagramSession  string // Instagram session cookie for authenticated requests
}

func Load() (*Config, error) {
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	notificationDelay, _ := time.ParseDuration(getEnv("NOTIFICATION_DELAY_HOURS", "6") + "h")

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/unfollow_tracker?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-me-in-production"),
			Expiry: jwtExpiry,
		},
		Scraper: ScraperConfig{
			DelayMs:           getEnvInt("SCRAPE_DELAY_MS", 2000),
			MaxConcurrent:     getEnvInt("MAX_CONCURRENT_FETCHERS", 4),
			ProxyPoolURL:      getEnv("PROXY_POOL_URL", ""),
			NotificationDelay: notificationDelay,
			InstagramSession:  getEnv("INSTAGRAM_SESSION_COOKIE", ""),
		},
	}, nil
}

// GetDatabaseConfig returns the database configuration for connection pooling
func (c *Config) GetDatabaseConfig() database.Config {
	return database.Config{
		URL:             c.Database.URL,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
