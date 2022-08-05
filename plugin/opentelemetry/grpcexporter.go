package opentelemetry

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func NewOtlpTraceGrpc(ctx context.Context, endpoint string, timeout time.Duration, header map[string]string) (*otlptrace.Exporter, error) {

	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	driver := otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: 2 * time.Microsecond,
			MaxInterval:     5 * time.Microsecond,
			// Never stop retry of retry-able status.
			MaxElapsedTime: 0,
		}),
		otlptracegrpc.WithTimeout(timeout*time.Second),
		otlptracegrpc.WithHeaders(header),
	)

	return otlptrace.New(ctx, driver)
}
