package ctxlog

import (
	"context"

	"go.uber.org/zap"
)

type loggerCtxKeyType struct{}

var loggerCtxKey = &loggerCtxKeyType{}

func AddZapField(ctx context.Context, field zap.Field) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	options, ok := ctx.Value(loggerCtxKey).([]zap.Field)
	if !ok {
		options = []zap.Field{}
	}

	options = append(options, field)
	return context.WithValue(ctx, loggerCtxKey, options)
}

func WithCtxData(ctx context.Context, log *zap.Logger) *zap.Logger {
	if ctx == nil {
		return log
	}

	options, ok := ctx.Value(loggerCtxKey).([]zap.Field)
	if ok {
		log = log.With(options...)
	}

	return log
}
