package config

import (
	"log/slog"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	DB   DBConfig
	Auth AuthConfig
}

type DBConfig struct {
	User     string `env:"DB_USER" envDefault:"root"`
	Password string `env:"DB_PASS" envDefault:""`
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"3306"`
	Database string `env:"DB_NAME" envDefault:"cloak"`
}

type AuthConfig struct {
	JWTSecret string `env:"AUTH_JWT_SECRET" envDefault:"secret"`
}

func Load() Config {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		slog.Error("Error loading config", err)
	}
	return config
}
