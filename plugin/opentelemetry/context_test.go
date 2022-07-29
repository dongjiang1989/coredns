package opentelemetry

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func testContext(provider trace.TracerProvider) context.Context {
	return withProvider(context.Background(), provider, propagation.TraceContext{}, "test")
}

func TestContext(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	ctx, err := Context(context.Background(), "test", WithBatchTimeout(3*time.Second), WithIsDropOnQueueFull(true),
		WithMaxExportBatchSize(256), WithMaxQueueSize(1024), WithExporter(exporter),
		WithResource(&resource.Resource{}), WithSamplingStrategy(AlwaysOn), WithInactiveTimeout(1*time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	s := ctx.Value(stateKey)
	if s == nil {
		t.Fatal("expected state in context")
	}
	st, ok := s.(*stateBag)
	if !ok {
		t.Fatalf("got %T, expected *stateBag", s)
	}
	if st.provider == nil {
		t.Error("expected provider in tracing context")
	}
}

func TestContextWith(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	ctx, err := Context(context.Background(), "test", WithBatchTimeout(3*time.Second), WithIsDropOnQueueFull(false),
		WithMaxExportBatchSize(256), WithMaxQueueSize(1024), WithExporter(exporter),
		WithResource(&resource.Resource{}), WithSamplingStrategy(AlwaysOn), WithInactiveTimeout(1*time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	s := ctx.Value(stateKey)
	if s == nil {
		t.Fatal("expected state in context")
	}
	st, ok := s.(*stateBag)
	if !ok {
		t.Fatalf("got %T, expected *stateBag", s)
	}
	if st.provider == nil {
		t.Error("expected provider in tracing context")
	}
}

func TestWithTracing(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	gctx := context.Background()
	ctx, err := Context(gctx, "test", WithBatchTimeout(3*time.Second), WithIsDropOnQueueFull(false),
		WithMaxExportBatchSize(256), WithMaxQueueSize(1024), WithExporter(exporter),
		WithResource(&resource.Resource{}), WithSamplingStrategy(AlwaysOn), WithInactiveTimeout(1*time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	s := ctx.Value(stateKey)
	if s == nil {
		t.Fatal("expected state in context")
	}
	st, ok := s.(*stateBag)
	if !ok {
		t.Fatalf("got %T, expected *stateBag", s)
	}
	if st.provider == nil {
		t.Error("expected provider in tracing context")
	}

	newctx := withTracing(ctx, gctx)
	if newctx == nil {
		t.Fatalf("got %v, expected cxt", newctx)
	}
}
