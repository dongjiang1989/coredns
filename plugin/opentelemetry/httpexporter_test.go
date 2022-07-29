package opentelemetry

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNewOtlpTraceHttp(t *testing.T) {

	ctx := context.Background()
	exporter, err := NewOtlpTraceHttp(ctx, "http://localhost:4318", time.Second, nil)

	if err != nil {
		t.Errorf("got %v", err)
	}
	defer func() {
		err = exporter.Shutdown(ctx)
		if err != nil {
			t.Errorf("got %v", err)
		}
	}()
	err = exporter.ExportSpans(ctx, nil)
	if os.IsTimeout(err) == true {
		t.Errorf("expected timeout error, got: %v", err)
	}
}
