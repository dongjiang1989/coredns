package opentelemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func NewOtlpTraceHttp(ctx context.Context, endpoint string, timeout time.Duration, header map[string]string) (*otlptrace.Exporter, error) {
	driver := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 2 * time.Microsecond,
			MaxInterval:     5 * time.Microsecond,
			// Never stop retry of retry-able status.
			MaxElapsedTime: 0,
		}),
		otlptracehttp.WithTimeout(timeout*time.Second),
		otlptracehttp.WithHeaders(header),
	)

	return otlptrace.New(ctx, driver)
}
