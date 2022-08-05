# opentelemetry

## Name

*opentelemetry* - enables opentelemetry trace of DNS requests as they go through the plugin chain. The Plugin supports binary-encoded OLTP over HTTP/GRPC.

## Description

With *opentelemetry* you enable opentelemetry of how a request flows through CoreDNS. Enable the *debug*
plugin to get logs from the trace plugin.

## Syntax

The simplest form is just:

~~~
opentelemetry {
	endpoint URLPATH
	sample NAME
}
~~~

* **URLPATH** is the type of tracing oltp server endpoint. Currently only `http://xxx` and `grpc://xxx` are supported.
* **NAME** is the sample strategy. Currently only `alwaysoff` , `traceidratio` and `alwayson` are supported. Default to `alwayson`

With this form, all queries will be traced.

Additional features can be enabled with this syntax:
		
~~~
opentelemetry [servicename] {
	sample NAME
	fraction RATE
	
	endpoint URLPATH
	endpoint_timeout DURATION
	endpoint_header ARRAY
	
	drop_on_queue_full BOOL
	max_queue_size SIZE
	batch_timeout DURATION
	max_export_batch_size SIZE
	inactive_timeout DURATION
}
~~~

* `servicename` is opentelemetry report service name. default "coredns".
* `sample` **NAME** is the sample strategy. Currently only `alwaysoff` , `traceidratio` and `alwayson` are supported. Default to `alwayson`.  The default is alwayson.
* `fraction` **RATE** is the ratio for sample strategy. Must in [0,1]. if sample is `alwayson`, fraction is `1`; if sample is `alwaysoff`,  fraction is `1`; if sample is `traceidratio`, using fraction value.

* `endpoint` **URLPATH** is the type of tracing OLTP server endpoint. Currently only `http://xxx` and `grpc://xxx` are supported.
* `endpoint_timeout` **DURATION** is the request timeout  of tracing OLTP server endpoint. default 3s.
* `endpoint_header` **ARRAY** is the request header  of tracing OLTP server endpoint. default "". Must "key1,value1,key2,value2,..."

* `drop_on_queue_full` **BOOL** drop span when queue is full, otherwise force process batches. default true
* `max_queue_size` **SIZE**  maximum queue size to buffer spans for delayed processing. default 2048
* `batch_timeout` **DURATION**  maximum duration(second) for constructing a batch. default 5s
* `max_export_batch_size` **SIZE**  maximum number of spans to process in a single batch. default 256
* `inactive_timeout` **DURATION**  timer interval(second) for processing batches. default 2s


Note the opentelemetry provider does not support the v1 API since coredns 1.7.1.

## Examples

Use an alternative OLTP address:

~~~
opentelemetry {
	endpoint http://tracinghost:4318
	sample alwayson
}
~~~

or

~~~ corefile
. {
    opentelemetry {
		endpoint http://tracinghost:4318
		sample alwayson
	}
}
~~~


## See Also

See the *debug* plugin for more information about debug logging.
