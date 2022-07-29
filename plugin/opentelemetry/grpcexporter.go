package opentelemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func NewOtlpTraceGrpc(ctx context.Context, endpoint string, timeout time.Duration, header map[string]string) (*otlptrace.Exporter, error) {
	driver := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(endpoint),
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
