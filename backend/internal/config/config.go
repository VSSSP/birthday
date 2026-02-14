package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Google   GoogleConfig
	Apple    AppleConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	Env  string `env:"SERVER_ENV" envDefault:"development"`
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	URL string `env:"DATABASE_URL,required"`
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	URL string `env:"REDIS_URL" envDefault:"redis://localhost:6379"`
}

// JWTConfig holds JWT signing settings.
type JWTConfig struct {
	AccessSecret  string        `env:"JWT_ACCESS_SECRET,required"`
	RefreshSecret string        `env:"JWT_REFRESH_SECRET,required"`
	AccessExpiry  time.Duration `env:"JWT_ACCESS_EXPIRY" envDefault:"15m"`
	RefreshExpiry time.Duration `env:"JWT_REFRESH_EXPIRY" envDefault:"168h"`
}

// GoogleConfig holds Google OAuth settings.
type GoogleConfig struct {
	ClientID string `env:"GOOGLE_CLIENT_ID" envDefault:""`
}

// AppleConfig holds Apple Sign In settings.
type AppleConfig struct {
	ClientID string `env:"APPLE_CLIENT_ID" envDefault:""`
}

// Load parses environment variables into a Config struct.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
