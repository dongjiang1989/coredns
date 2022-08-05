package opentelemetry

import (
	"context"
	"fmt"
	"time"

	"errors"
	//"net/http/httptest"
	"testing"

	"github.com/coredns/caddy"
	//"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"

	//"github.com/coredns/coredns/plugin/pkg/rcode"
	"github.com/coredns/coredns/plugin/test"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	//"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestNewTelemetry(t *testing.T) {
	m, err := parse(caddy.NewTestController("dns", `opentelemetry coredns { 
	sample ratio
	fraction 0.5
	
    endpoint http://localhost:8080 
	endpoint_timeout 2s
	endpoint_header key1=val1 key2=val2
	
	drop_on_queue_full
	max_queue_size 2000
	batch_timeout 1s
	max_export_batch_size 128
	inactive_timeout 1s
}`))
	if err != nil {
		t.Errorf("Error parsing test input: %s", err)
		return
	}

	if m.Name() != "opentelemetry" {
		t.Errorf("Wrong name from GetName: %s", m.Name())
	}

	if len(m.opts.collectorRequestHeader) != 2 {
		t.Errorf("Wrong collectorRequestHeader from GetName: %s", m.opts.collectorRequestHeader)
	}

	if m.opts.collectorRequestTimeout != 2*time.Second {
		t.Errorf("Wrong collectorRequestTimeout from collectorRequestTimeout: %s", m.opts.collectorRequestTimeout)
	}

	if m.opts.maxQueueSize != 2000 {
		t.Errorf("Wrong max_queue_size from maxQueueSize: %v", m.opts.maxQueueSize)
	}

	if m.opts.isDropOnQueueFull != true {
		t.Errorf("Wrong drop_on_queue_full from isDropOnQueueFull: %v", m.opts.isDropOnQueueFull)
	}

	if m.opts.collectorEndpoint != "http://localhost:8080" {
		t.Errorf("Wrong collectorEndpoint from collectorEndpoint: %s", m.opts.collectorEndpoint)
	}

	err = m.NewTelemetry(context.Background())
	if err != nil {
		t.Errorf("Error starting opentelemetry plugin: %s", err)
		return
	}

	if m.serviceName != "coredns" {
		t.Errorf("serviceName is ''")
	}

	if m.TracerProvider() == nil {
		t.Errorf("Error, no tracer created")
	}
}
