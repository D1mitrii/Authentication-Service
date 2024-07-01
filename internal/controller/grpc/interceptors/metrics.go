package interceptors

import (
	"context"
	"github.com/d1mitrii/authentication-service/internal/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MetricsInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	diff := time.Since(start).Seconds()
	status := status.Code(err).String()

	metrics.GrpcCounterRequestTotal(status, info.FullMethod)
	metrics.GrpcHistogramResponseTimeObserve(status, info.FullMethod, diff)

	return resp, err
}
