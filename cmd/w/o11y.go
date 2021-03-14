package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

func o11y(ctx context.Context) (func(context.Context), http.Handler, error) {
	name, version := "go.seankhliao.com/w/cmd/w", "v15"
	bi, ok := debug.ReadBuildInfo()
	if ok {
		name = bi.Main.Path
		version = bi.Main.Version
	}

	res := resource.NewWithAttributes(
		semconv.ServiceNameKey.String(name),
		semconv.ServiceVersionKey.String(version),
	)

	// trace and metrics exporter
	// exporter, err := otlp.NewExporter(ctx, otlpgrpc.NewDriver(
	// 	otlpgrpc.WithInsecure(),
	// 	otlpgrpc.WithEndpoint("otel-collector.otel.svc.cluster.local:55680"),
	// ))
	// if err != nil {
	// 	return nil, fmt.Errorf("o11y: create exporter: %w", err)
	// }

	// metrics exporter
	hf, err := prometheus.InstallNewPipeline(
		prometheus.Config{},
		basic.WithResource(res),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("o11y: create prometheus: %w", err)
	}

	// default metrics
	err = runtime.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("o11y: start runtime metrics: %w", err)
	}
	err = host.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("o11y: start host metrics: %w", err)
	}

	// trace provider
	traceProvider := trace.NewTracerProvider(trace.WithConfig(
		trace.Config{DefaultSampler: trace.AlwaysSample()}),
	// sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	))

	return func(ctx context.Context) {
			err := traceProvider.Shutdown(ctx)
			if err != nil {
				otel.Handle(err)
			}
			// err = exporter.Shutdown(ctx)
			// if err != nil {
			// 	otel.Handle(err)
			// }
		},
		hf,
		nil
}
