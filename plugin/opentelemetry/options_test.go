package opentelemetry

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

/*
func defaultOptions() *options {
	return &options{
		samplingStrategy:         AlwaysOn,
		samplingStrategyFraction: 1.0,
		isAddMetadata:            true,
		propagator:               propagation.TraceContext{},
		isDropOnQueueFull:        true,
		maxQueueSize:             2048,
		batchTimeout:             5,
		maxExportBatchSize:       256,
		inactiveTimeout:          2,
		collectorAddr:            "localhost:9445",
		collectorRequestTimeout:  3,
	}
}
*/
func TestOptions(t *testing.T) {
	ctx := context.Background()

	options := defaultOptions()
	if options.samplingStrategy != AlwaysOn {
		t.Errorf("got %v, want 1", options.samplingStrategy)
	}
	if options.samplingStrategyFraction != 1.0 {
		t.Errorf("got %v, want 1.0", options.samplingStrategyFraction)
	}
	if options.isAddMetadata != true {
		t.Errorf("got %v, want true", options.isAddMetadata)
	}
	if options.isDropOnQueueFull != true {
		t.Errorf("got %v, want true", options.isDropOnQueueFull)
	}
	if options.maxQueueSize != 2048 {
		t.Errorf("got %v, want 2048", options.maxQueueSize)
	}
	if options.batchTimeout != 5*time.Second {
		t.Errorf("got %v, want 5", options.batchTimeout)
	}
	if options.maxExportBatchSize != 256 {
		t.Errorf("got %v, want 256", options.maxExportBatchSize)
	}
	if options.inactiveTimeout != 2*time.Second {
		t.Errorf("got %v, want 2", options.inactiveTimeout)
	}
	if options.collectorEndpoint != "" {
		t.Errorf("got %v, want ''", options.collectorEndpoint)
	}
	if options.collectorRequestTimeout != 3*time.Second {
		t.Errorf("got %v, want 3", options.collectorRequestTimeout)
	}

	WithSamplingStrategy(AlwaysOff)(ctx, options)
	if options.samplingStrategy != AlwaysOff {
		t.Errorf("got %v sampling rate, want 0", options.samplingStrategy)
	}
	if options.samplingStrategyFraction != 0.0 {
		t.Errorf("got %v, want 0.0", options.samplingStrategyFraction)
	}

	WithExporter(tracetest.NewInMemoryExporter())(ctx, options)
	if options.exporter == nil {
		t.Error("got nil exporter, want non-nil")
	}

	WithResource(&resource.Resource{})(ctx, options)
	if options.resource == nil {
		t.Error("got nil resource, want non-nil")
	}

	WithSamplingStrategy(TraceIdRatio)(ctx, options)
	WithSamplingStrategyFraction(0.2)(ctx, options)
	if options.samplingStrategy != TraceIdRatio {
		t.Errorf("got %v, want 2", options.samplingStrategy)
	}
	if options.samplingStrategyFraction != 0.2 {
		t.Errorf("got %v, want 0.2", options.samplingStrategyFraction)
	}

	WithIsAddMetadata(false)(ctx, options)
	if options.isAddMetadata != false {
		t.Errorf("got %v, want false", options.isAddMetadata)
	}

	WithIsDropOnQueueFull(false)(ctx, options)
	if options.isDropOnQueueFull != false {
		t.Errorf("got %v, want false", options.isDropOnQueueFull)
	}

	WithMaxQueueSize(1000)(ctx, options)
	if options.maxQueueSize != 1000 {
		t.Errorf("got %v, want 1000", options.maxQueueSize)
	}

	WithBatchTimeout(3*time.Second)(ctx, options)
	if options.batchTimeout != 3*time.Second {
		t.Errorf("got %v, want 3", options.batchTimeout)
	}

	WithMaxExportBatchSize(128)(ctx, options)
	if options.maxExportBatchSize != 128 {
		t.Errorf("got %v, want 128", options.maxExportBatchSize)
	}

	WithInactiveTimeout(1*time.Second)(ctx, options)
	if options.inactiveTimeout != 1*time.Second {
		t.Errorf("got %v, want 1", options.inactiveTimeout)
	}

	WithCollectorEndpoint("mock address")(ctx, options)
	if options.collectorEndpoint != "mock address" {
		t.Errorf("got %v, want mock address", options.collectorEndpoint)
	}

	WithCollectorRequestTimeout(1*time.Second)(ctx, options)
	if options.collectorRequestTimeout != 1*time.Second {
		t.Errorf("got %v, want 1", options.collectorRequestTimeout)
	}

	WithAttributes([]attribute.KeyValue{attribute.String("k1", "v11"), attribute.String("k1", "v12"), attribute.String("k2", "v21")})(ctx, options)
	if len(options.resource.Attributes()) != len([]attribute.KeyValue{attribute.String("k1", "v11"), attribute.String("k2", "v21")}) {
		t.Error("got  resource, want non-nil")
	}
	if options.attributes == nil {
		t.Error("got nil attributes, want non-nil")
	}
}
