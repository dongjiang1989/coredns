package opentelemetry

import (
	"context"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type (
	options struct {
		samplingStrategy         SamplingStrategy
		samplingStrategyFraction float64 // trace_id_ratio fraction [0,1]

		isAddMetadata bool // add metadata plugin data

		attributes []attribute.KeyValue

		collectorEndpoint       string
		collectorSchema         CollectorSchema
		collectorRequestTimeout time.Duration // default 3s
		collectorRequestHeader  map[string]string

		exporter   sdktrace.SpanExporter
		propagator propagation.TextMapPropagator
		//parentSamplerOptions []sdktrace.ParentBasedSamplerOption
		resource           *resource.Resource
		isDropOnQueueFull  bool          // drop span when queue is full, otherwise force process batches. default true
		maxQueueSize       int           // maximum queue size to buffer spans for delayed processing. default 2048
		batchTimeout       time.Duration // maximum duration(second) for constructing a batch. default 5s
		maxExportBatchSize int           // maximum number of spans to process in a single batch. default 256
		inactiveTimeout    time.Duration // timer interval(second) for processing batches. default 2s
	}

	// TraceOption is a function that configures a provider.
	TraceOption func(ctx context.Context, opts *options) error
)

// defaultOptions returns the default sampler options.
func defaultOptions() *options {
	return &options{
		samplingStrategy:         AlwaysOn,
		samplingStrategyFraction: 1.0,
		isAddMetadata:            true,
		propagator:               propagation.TraceContext{},
		isDropOnQueueFull:        true,
		maxQueueSize:             2048,
		batchTimeout:             5 * time.Second,
		maxExportBatchSize:       256,
		inactiveTimeout:          2 * time.Second,
		collectorRequestTimeout:  3 * time.Second,
	}
}

// WithSamplingStrategy sets sampling strategy.
func WithSamplingStrategy(sampling SamplingStrategy) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.samplingStrategy = sampling
		if opts.samplingStrategy == AlwaysOn {
			opts.samplingStrategyFraction = 1.0
		} else if opts.samplingStrategy == AlwaysOff {
			opts.samplingStrategyFraction = 0.0
		}
		return nil
	}
}

// WithSamplingStrategyFraction sets the maximum sampling rate.
func WithSamplingStrategyFraction(rate float64) TraceOption {
	return func(ctx context.Context, opts *options) error {
		if opts.samplingStrategy == AlwaysOn {
			opts.samplingStrategyFraction = 1.0
		} else if opts.samplingStrategy == AlwaysOff {
			opts.samplingStrategyFraction = 0.0
		} else {
			opts.samplingStrategyFraction = rate
		}
		return nil
	}
}

// WithIsAddMetadata sets add metadate to trace.
func WithIsAddMetadata(ok bool) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.isAddMetadata = ok
		return nil
	}
}

// WithExporter sets the exporter to use.
func WithExporter(exporter sdktrace.SpanExporter) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.exporter = exporter
		return nil
	}
}

// WithCollectorEndpoint sets the address for collctor.
func WithCollectorEndpoint(endpoint string) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.collectorEndpoint = endpoint
		return nil
	}
}

// WithCollectorRequestTimeout sets maximum duration(second) for collctor.
func WithCollectorRequestTimeout(timeout time.Duration) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.collectorRequestTimeout = timeout
		return nil
	}
}

// WithResource sets the underlying opentelemetry resource.
func WithResource(res *resource.Resource) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.resource = res
		return nil
	}
}

// WithIsDropOnQueueFull sets drop span when queue is full
func WithIsDropOnQueueFull(drop bool) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.isDropOnQueueFull = drop
		return nil
	}
}

// WithMaxQueueSize sets maximum queue size
func WithMaxQueueSize(size int) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.maxQueueSize = size
		return nil
	}
}

// WithBatchTimeout sets maximum duration(second) for constructing a batch
func WithBatchTimeout(timeout time.Duration) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.batchTimeout = timeout
		return nil
	}
}

// WithMaxExportBatchSize sets maximum number of spans to process in a single batch
func WithMaxExportBatchSize(size int) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.maxExportBatchSize = size
		return nil
	}
}

// WithInactiveTimeout sets timer interval(second) for processing batches
func WithInactiveTimeout(timeout time.Duration) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.inactiveTimeout = timeout
		return nil
	}
}

//WithAttributes set attrs custom kv message
func WithAttributes(attrs []attribute.KeyValue) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.attributes = attrs
		opts.resource = resource.NewSchemaless(attrs...)
		return nil
	}
}

//WithCustomExporter set user custom exporter
func WithCustomExporter(endpoint string, timeout time.Duration, header map[string]string) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.collectorEndpoint = endpoint
		opts.collectorRequestTimeout = timeout
		opts.collectorRequestHeader = header

		u, err := url.Parse(endpoint)
		if err != nil {
			return errCollectorEndpointStyle
		}
		switch strings.ToLower(u.Scheme) {
		case "http":
			opts.exporter, err = NewOtlpTraceHttp(ctx, endpoint, timeout, header)
			if err != nil {
				return err
			}
			return nil
		case "grpc":
			opts.exporter, err = NewOtlpTraceGrpc(ctx, endpoint, timeout, header)
			if err != nil {
				return err
			}
			return nil
		default:
			return errCollectorEndpointStyle
		}

		return nil
	}
}
