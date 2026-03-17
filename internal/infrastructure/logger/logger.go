package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Environment string

const (
	EnvProduction  Environment = "production"
	EnvDevelopment Environment = "development"
)

func New(env Environment) (*zap.Logger, error) {
	switch env {
	case EnvProduction:
		return zap.NewProduction()
	case EnvDevelopment:
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return cfg.Build()
	default:
		return nil, fmt.Errorf("unknown environment: %s", env)
	}
}
