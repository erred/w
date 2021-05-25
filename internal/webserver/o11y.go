package webserver

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"google.golang.org/grpc"
)

func o11y(ctx context.Context, endpoint string) (*otlp.Exporter, http.Handler, error) {
	r := newResource()
	exp, err := newTrace(ctx, r, endpoint)
	if err != nil {
		return nil, nil, err
	}
	h, err := newMetrics(r)
	if err != nil {
		return nil, nil, err
	}

	return exp, h, err
}

func newResource() *resource.Resource {
	name, version := "go.seankhliao.com/w/cmd/w", "v15"
	bi, ok := debug.ReadBuildInfo()
	if ok {
		name = bi.Main.Path
		version = bi.Main.Version
	}

	return resource.NewWithAttributes(
		semconv.ServiceNameKey.String(name),
		semconv.ServiceVersionKey.String(version),
	)
}

func newTrace(ctx context.Context, r *resource.Resource, endpoint string) (*otlp.Exporter, error) {
	if endpoint == "" {
		return nil, nil
	}

	// exporter
	exporter, err := otlp.NewExporter(
		ctx,
		otlpgrpc.NewDriver(
			otlpgrpc.WithInsecure(),
			otlpgrpc.WithEndpoint(endpoint),
			otlpgrpc.WithDialOption(
				grpc.WithBlock(),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	// trace provider
	traceProvider := trace.NewTracerProvider(
		trace.WithResource(r),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
	)

	// register
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	))
	return exporter, nil
}

func newMetrics(r *resource.Resource) (http.Handler, error) {
	h, err := prometheus.InstallNewPipeline(
		prometheus.Config{},
		basic.WithResource(r),
	)
	if err != nil {
		return nil, fmt.Errorf("o11y: create prometheus: %w", err)
	}

	// default metrics
	err = runtime.Start()
	if err != nil {
		return nil, fmt.Errorf("o11y: start runtime metrics: %w", err)
	}
	err = host.Start()
	if err != nil {
		return nil, fmt.Errorf("o11y: start host metrics: %w", err)
	}

	return h, nil
}
