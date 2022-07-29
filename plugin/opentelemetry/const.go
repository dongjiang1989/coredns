package opentelemetry

import (
	"errors"
)

const (

	// CoreDNS report trace serive name
	defaultServiceName = "coredns/coredns"

	// CoreDNS opentelemetry plugin register name
	pluginName = "opentelemetry"

	// AttributeRequestID is the name of the span attribute that contains the
	// request ID.
	attributeRequestID = "request.id"

	//Metadata plugin traceid key
	metaTraceIdKey = "telemetry/traceid"
)

// error
var (
	errCollectorEndpointStyle = errors.New("OpenTelemetry: Collector Endpoint Style Error: Must http:// or grpc:// start")
)

type SamplingStrategy int

const (
	AlwaysOff    SamplingStrategy = iota // sampling nothing
	AlwaysOn                             // sampling all
	TraceIdRatio                         // base trace id percentage
)

type TracingExporterType int

const (
	Tracing TracingExporterType = iota
	Metrics
	Logging
)

type CollectorSchema int

const (
	Http CollectorSchema = iota + 1
	Grpc
)
