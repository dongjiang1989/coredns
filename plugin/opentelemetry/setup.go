package opentelemetry

import (
	//"context"
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin(pluginName)

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	t, err := parse(c)
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		t.Next = next
		return t
	})

	//t.NewTelemetry(context.Background())
	//c.OnStartup(t.NewTelemetry(c.Context()))
	//c.OnStartup(t.OnStartup)

	return nil
}

func parse(c *caddy.Controller) (*telemetry, error) {
	var err error
	var serviceName string
	opts := defaultOptions()

	cfg := dnsserver.GetConfig(c)
	if cfg.ListenHosts[0] != "" {
		opts.collectorEndpoint = cfg.ListenHosts[0] + ":" + cfg.Port
	}
	/*
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
	*/
	for c.Next() { // telemetry
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			serviceName = defaultServiceName
		case 1:
			serviceName = args[0]
		default:
			err = c.ArgErr()
		}
		if err != nil {
			return nil, err
		}
		for c.NextBlock() {
			switch c.Val() {
			case "sample":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				opts.samplingStrategy = Value(strings.ToLower(args[0]))
			case "fraction":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				opts.samplingStrategyFraction, err = strconv.ParseFloat(args[0], 64)
				if err != nil {
					return nil, err
				}
				if opts.samplingStrategyFraction > 1 || opts.samplingStrategyFraction < 0 {
					return nil, fmt.Errorf("opentelemetry sampling strategy fraction rate must be between 0 and 1, '%f' is not supported", opts.samplingStrategyFraction)
				}
			case "endpoint":
				args := c.RemainingArgs()
				if len(args) > 1 {
					return nil, c.ArgErr()
				}
				if len(args) == 1 {
					opts.collectorEndpoint = args[0]
				}
			case "endpoint_timeout":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				opts.collectorRequestTimeout, err = time.ParseDuration(args[0])
				if err != nil {
					return nil, err
				}
			case "endpoint_header":
				args := c.RemainingArgs()
				if len(args) < 1 {
					return nil, c.ArgErr()
				}
				opts.collectorRequestHeader, err = ArrayToMap(args)
				if err != nil {
					return nil, err
				}
			case "drop_on_queue_full":
				args := c.RemainingArgs()
				if len(args) > 1 {
					return nil, c.ArgErr()
				} else if len(args) == 1 {
					opts.isDropOnQueueFull, err = strconv.ParseBool(args[0])
					if err != nil {
						opts.isDropOnQueueFull = true
					}
				} else {
					opts.isDropOnQueueFull = true
				}
			case "max_queue_size":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				opts.maxQueueSize, err = strconv.Atoi(args[0])
				if err != nil {
					return nil, err
				}
			case "batch_timeout":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				opts.batchTimeout, err = time.ParseDuration(args[0])
				if err != nil {
					return nil, err
				}
			case "max_export_batch_size":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				opts.maxExportBatchSize, err = strconv.Atoi(args[0])
				if err != nil {
					return nil, err
				}
			case "inactive_timeout":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				opts.inactiveTimeout, err = time.ParseDuration(args[0])
				if err != nil {
					return nil, err
				}
			}
		}
	}
	clog.Debug("opentelemetry options: ", opts)
	return &telemetry{
		opts:        opts,
		serviceName: serviceName,
	}, err
}

func ArrayToMap(arr []string) (map[string]string, error) {
	res := make(map[string]string)
	for _, str := range arr {
		kv := strings.Split(str, "=")
		if len(kv) != 2 {
			return nil, errCollectorRequestHeaderStyle
		}
		res[kv[0]] = kv[1]
	}
	return res, nil
}
