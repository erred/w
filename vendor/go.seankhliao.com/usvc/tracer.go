package usvc

import (
	"flag"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type TracerOpts struct {
	Enabled           bool
	CollectorEndpoint string
}

func (o *TracerOpts) Flags(fs *flag.FlagSet) {
	fs.BoolVar(&o.Enabled, "trace", false, "enable tracing")
	fs.StringVar(&o.CollectorEndpoint, "trace.collector", "http://jaeger:14268/api/traces?format=jaeger.thrift", "jaeger collector endpoint")
}

// Tracer installs a global tracer
func (o TracerOpts) Tracer() (shutdown func() error, err error) {
	var opts []jaeger.Option
	opts = append(opts, jaeger.WithProcess(jaeger.Process{
		ServiceName: name,
	}))
	opts = append(opts, jaeger.WithProcessFromEnv())
	if !o.Enabled {
		opts = append(opts, jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.NeverSample()}))
	}

	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(
			o.CollectorEndpoint,
			jaeger.WithCollectorEndpointOptionFromEnv(),
		),
		opts...,
	)
	shutdown = func() error {
		flush()
		return nil
	}

	b3 := b3.B3{}
	global.SetPropagators(propagation.New(
		propagation.WithExtractors(b3),
		propagation.WithInjectors(b3),
	))
	return
}
