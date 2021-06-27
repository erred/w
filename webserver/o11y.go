package webserver

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type stopper interface {
	Stop(ctx context.Context) error
}

type shutdowner interface {
	Shutdown(ctx context.Context) error
}

type shutdown struct {
	s1 []shutdowner
	s2 []stopper
}

func (s *shutdown) Shutdown(ctx context.Context) error {
	for _, sd := range s.s1 {
		otel.Handle(sd.Shutdown(ctx))
	}
	for _, sd := range s.s2 {
		otel.Handle(sd.Stop(ctx))
	}
	return nil
}

func o11y(ctx context.Context, endpoint string) (*shutdown, error) {
	if endpoint == "" {
		return nil, nil
	}

	// trace
	traceExporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	))
	if err != nil {
		return nil, fmt.Errorf("create otlptrace exporter: %w", err)
	}
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
		jaeger.Jaeger{},
	))

	// metric
	metricExporter, err := otlpmetric.New(ctx, otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint),
	))
	if err != nil {
		return nil, fmt.Errorf("create otlpmetric exporter: %w", err)
	}
	pusher := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			metricExporter,
		),
		controller.WithExporter(metricExporter),
		controller.WithCollectPeriod(2*time.Second),
	)
	global.SetMeterProvider(pusher.MeterProvider())

	return &shutdown{[]shutdowner{traceExporter, traceProvider, metricExporter}, []stopper{pusher}}, nil
}
