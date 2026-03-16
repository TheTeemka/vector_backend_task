package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort string         `env:"GRPC_PORT" validate:"required"`
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST"     validate:"required"`
	Port     string `env:"POSTGRES_PORT"     validate:"required"`
	User     string `env:"POSTGRES_USER"     validate:"required"`
	Password string `env:"POSTGRES_PASSWORD" validate:"required"`
	Name     string `env:"POSTGRES_DB"       validate:"required"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" validate:"required,oneof=disable require verify-ca verify-full"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if err := validator.New().Struct(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
