package opentelemetry

import (
	"context"
	//"fmt"
	"net/http"
	"sync"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/pkg/rcode"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"

	"go.opentelemetry.io/otel"
	//"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type telemetry struct {
	Next        plugin.Handler
	serviceName string
	opts        *options
	ctx         context.Context
	Once        sync.Once
}

// NewTelemetry sets up the telemetry
func (t *telemetry) NewTelemetry(ctx context.Context) error {
	var err error
	var gctx context.Context
	t.Once.Do(func() {
		gctx, err = Context(ctx, t.serviceName,
			WithAttributes(t.opts.attributes),
			WithBatchTimeout(t.opts.batchTimeout),
			WithCustomExporter(t.opts.collectorEndpoint, t.opts.collectorRequestTimeout, t.opts.collectorRequestHeader),
			WithInactiveTimeout(t.opts.inactiveTimeout),
			WithIsAddMetadata(t.opts.isAddMetadata),
			WithIsDropOnQueueFull(t.opts.isDropOnQueueFull),
			WithMaxExportBatchSize(t.opts.maxExportBatchSize),
			WithMaxQueueSize(t.opts.maxQueueSize),
			WithSamplingStrategy(t.opts.samplingStrategy),
			WithSamplingStrategyFraction(t.opts.samplingStrategyFraction))
		t.ctx = gctx
	})
	return err
}

// TracerProvider implements
func (t *telemetry) TracerProvider() trace.TracerProvider {
	return TraceProvider(t.ctx)
}

// Name implements the Handler interface.
func (t *telemetry) Name() string { return pluginName }

// ServeDNS implements the plugin.Handle interface.
func (t *telemetry) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	ctx = StartSpan(ctx, defaultServiceName)
	if trace.SpanFromContext(ctx).SpanContext().IsValid() {
		return plugin.NextOrFailure(t.Name(), t.Next, ctx, w, r)
	}

	if val := ctx.Value(dnsserver.HTTPRequestKey{}); val != nil {
		if httpReq, ok := val.(*http.Request); ok {
			//t.t.Extract(ctx, propagation.HeaderCarrier(httpReq.Header))
			//SetSpanAttributes(ctx, propagation.HeaderCarrier(httpReq.Header))
			otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(httpReq.Header))
		}
	}

	req := request.Request{W: w, Req: r}
	defer EndSpan(ctx)

	metadata.SetValueFunc(ctx, metaTraceIdKey, func() string { return TraceID(ctx) })

	rw := dnstest.NewRecorder(w)
	status, err := plugin.NextOrFailure(t.Name(), t.Next, ctx, rw, r)

	rc := rw.Rcode
	if !plugin.ClientWrite(status) {
		rc = status
	}

	SetSpanAttributes(ctx, Name, req.Name(),
		Type, req.Type(),
		Proto, req.Proto(),
		Remote, req.IP(),
		Rcode, rcode.ToString(rc))

	return status, err
}
