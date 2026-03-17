package interceptor

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	zaplogger "shipment-service/internal/infrastructure/logger"
)

func UnaryLogger(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		reqID := uuid.NewString()
		ctx = zaplogger.AddZapField(ctx, zap.String("request_id", reqID))

		start := time.Now()
		resp, err := handler(ctx, req)

		zaplogger.WithCtxData(ctx, log).Info("grpc call",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)
		return resp, err
	}
}

func UnaryRecovery(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				zaplogger.WithCtxData(ctx, log).Error("panic recovered",
					zap.Any("panic", r),
					zap.String("method", info.FullMethod),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
