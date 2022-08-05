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

const (
	Name   = "coredns.io/name"
	Type   = "coredns.io/type"
	Rcode  = "coredns.io/rcode"
	Proto  = "coredns.io/proto"
	Remote = "coredns.io/remote"
)

// error
var (
	errCollectorEndpointStyle      = errors.New("OpenTelemetry: Collector Endpoint Style Error: Must http:// or grpc:// start")
	errCollectorRequestHeaderStyle = errors.New("OpenTelemetry: Options is not key=value style")
)

type SamplingStrategy int

const (
	Unknown      SamplingStrategy = -1
	AlwaysOff    SamplingStrategy = iota // sampling nothing
	AlwaysOn                             // sampling all
	TraceIdRatio                         // base trace id percentage
)

func (s SamplingStrategy) String() string {
	if s == AlwaysOn {
		return "alwayson"
	} else if s == TraceIdRatio {
		return "ratio"
	} else {
		return "alwaysoff"
	}
}

func Value(name string) SamplingStrategy {
	if name == "alwaysoff" {
		return AlwaysOff
	} else if name == "alwayson" {
		return AlwaysOn
	} else if name == "ratio" {
		return TraceIdRatio
	} else {
		return Unknown
	}
}

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
