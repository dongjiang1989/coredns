package opentelemetry

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.8.0"
	"go.opentelemetry.io/otel/trace"
)

type (
	// ctxKey is a private type used to store the tracer provider in the context.
	ctxKey int

	// stateBag tracks the provider, tracer and active span sequence for a request.
	stateBag struct {
		svc        string
		provider   trace.TracerProvider
		propagator propagation.TextMapPropagator
		tracer     trace.Tracer
		spans      []trace.Span
	}
)

const (
	// stateKey is used to store the tracing state the context.
	stateKey ctxKey = iota + 1
)

// Context initializes the context so it can be used to create traces.
func Context(ctx context.Context, svc string, opts ...TraceOption) (context.Context, error) {
	options := defaultOptions()
	for _, o := range opts {
		err := o(ctx, options)
		if err != nil {
			return nil, err
		}
	}

	if options.exporter == nil {
		return nil, errors.New("missing exporter")
	}

	res := options.resource
	if res == nil {
		res = resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(svc))
	}

	var bsp sdktrace.SpanProcessor
	if options.isDropOnQueueFull {
		bsp = sdktrace.NewBatchSpanProcessor(options.exporter,
			sdktrace.WithBatchTimeout(options.batchTimeout),
			sdktrace.WithMaxExportBatchSize(options.maxExportBatchSize),
			sdktrace.WithMaxQueueSize(options.maxQueueSize),
			sdktrace.WithExportTimeout(options.batchTimeout),
			sdktrace.WithBlocking(),
		)
	} else {
		bsp = sdktrace.NewBatchSpanProcessor(options.exporter,
			sdktrace.WithBatchTimeout(options.batchTimeout),
			sdktrace.WithMaxExportBatchSize(options.maxExportBatchSize),
			sdktrace.WithMaxQueueSize(options.maxQueueSize),
			sdktrace.WithExportTimeout(options.batchTimeout),
		)
	}

	var provider *sdktrace.TracerProvider

	if options.samplingStrategy == AlwaysOff {
		provider = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.NeverSample()),
			sdktrace.WithResource(res),
			sdktrace.WithSpanProcessor(bsp),
		)
	} else if options.samplingStrategy == TraceIdRatio {
		provider = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.TraceIDRatioBased(options.samplingStrategyFraction)),
			sdktrace.WithResource(res),
			sdktrace.WithSpanProcessor(bsp),
		)
	} else {
		provider = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(res),
			sdktrace.WithSpanProcessor(bsp),
		)
	}

	return withProvider(ctx, provider, options.propagator, svc), nil
}

// TraceProvider returns the underlying otel trace provider.
func TraceProvider(ctx context.Context) trace.TracerProvider {
	sb := ctx.Value(stateKey).(*stateBag)
	return sb.provider
}

// withProvider stores the tracer provider in the context.
func withProvider(
	ctx context.Context,
	provider trace.TracerProvider,
	propagator propagation.TextMapPropagator,
	svc string) context.Context {

	return context.WithValue(
		ctx,
		stateKey,
		&stateBag{provider: provider, propagator: propagator, svc: svc},
	)
}

// withTracing initializes the tracing context, ctx must have been initialized
// with withProvider and the request must be traced by otel.
func withTracing(traceCtx, ctx context.Context) context.Context {
	state := traceCtx.Value(stateKey).(*stateBag)
	svc := state.svc
	provider := state.provider
	propagator := state.propagator
	tracer := provider.Tracer(defaultServiceName)
	spans := []trace.Span{trace.SpanFromContext(ctx)}
	return context.WithValue(ctx, stateKey, &stateBag{
		svc:        svc,
		provider:   provider,
		propagator: propagator,
		tracer:     tracer,
		spans:      spans,
	})
}
