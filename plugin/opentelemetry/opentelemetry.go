package opentelemetry

import (
	"context"
	//"fmt"
	//stdlog "log"
	//"net/http"
	"sync"
	//"sync/atomic"
	//"time"

	//"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	//"github.com/coredns/coredns/plugin/metadata"
	//"github.com/coredns/coredns/plugin/pkg/dnstest"
	//clog "github.com/coredns/coredns/plugin/pkg/log"
	//"github.com/coredns/coredns/plugin/pkg/rcode"
	_ "github.com/coredns/coredns/plugin/pkg/trace" // Plugin the trace package.
	//"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"go.opentelemetry.io/otel/attribute"
)

type telemetry struct {
	Next        plugin.Handler
	serviceName string
	exporter    TracingExporterType
	opt         options
	Once        sync.Once
	Attributes  []attribute.KeyValue
}

// OnStartup sets up the telemetry
func (t *telemetry) OnStartup() error {
	var err error
	t.Once.Do(func() {
		/*
			switch t.EndpointType {
			case "zipkin":
				err = t.setupZipkin()
			case "datadog":
				tracer := opentracer.New(
					tracer.WithAgentAddr(t.Endpoint),
					tracer.WithDebugMode(clog.D.Value()),
					tracer.WithGlobalTag(ext.SpanTypeDNS, true),
					tracer.WithServiceName(t.serviceName),
					tracer.WithAnalyticsRate(t.datadogAnalyticsRate),
					tracer.WithLogger(&loggerAdapter{log}),
				)
				t.tracer = tracer
				t.tagSet = tagByProvider["datadog"]
			default:
				err = fmt.Errorf("unknown endpoint type: %s", t.EndpointType)
			}
		*/
	})
	return err
}

// Name implements the Handler interface.
func (t *telemetry) Name() string { return pluginName }

// ServeDNS implements the plugin.Handle interface.
func (t *telemetry) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return 0, nil
}
