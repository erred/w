package usvc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/api/unit"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

type Components struct {
	Log    zerolog.Logger
	Meter  metric.Meter
	Tracer trace.Tracer
	HTTP   *http.ServeMux
	GRPC   *grpc.Server

	latency metric.Int64ValueRecorder

	flush func()
	prom  *prometheus.Exporter
}

func NewComponents(name, logLevel, logFormat string) (*Components, error) {
	c := &Components{}

	c.logger(logLevel, logFormat)

	err := c.tracer(name)
	if err != nil {
		return nil, fmt.Errorf("install tracer: %w", err)
	}

	err = c.metric(name)
	if err != nil {
		return nil, fmt.Errorf("install metric: %w", err)
	}

	return c, nil
}

func (c *Components) logger(logLevel, logFormat string) {
	var logout io.Writer = os.Stderr
	lvl, _ := zerolog.ParseLevel(logLevel)
	if logFormat == "logfmt" {
		logout = zerolog.ConsoleWriter{Out: logout}
	}
	c.Log = zerolog.New(logout).Level(lvl).With().Timestamp().Logger()
	log.SetOutput(c.Log)
}

func (c *Components) tracer(name string) error {
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("0.0.0.0:1234", jaeger.WithCollectorEndpointOptionFromEnv()),
		// jaeger.WithProcess(jaeger.Process{ServiceName: name}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.NeverSample()}),
		// jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	c.Tracer = global.Tracer(name)
	c.flush = flush

	// TODO: propagation

	return err
}

func (c *Components) metric(name string) error {
	promExporter, err := prometheus.InstallNewPipeline(prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 5, 10, 50, 100},
	})
	if err != nil {
		return err
	}
	c.prom = promExporter
	c.Meter = global.Meter(name)
	otelruntime.Start()
	c.latency = metric.Must(c.Meter).NewInt64ValueRecorder(
		"request_serve_latency_ms",
		metric.WithDescription("time to handle request"),
		metric.WithUnit(unit.Milliseconds),
	)
	return nil
}

func (c *Components) httpLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		defer func() {
			d := time.Since(t)
			c.latency.Record(r.Context(), d.Milliseconds())

			remote := r.Header.Get("x-forwarded-for")
			if remote == "" {
				remote = r.RemoteAddr
			}
			c.Log.Trace().
				Str("src", remote).
				Str("url", r.URL.String()).
				Str("user-agent", r.UserAgent()).
				Dur("dur", d).
				Msg("served")
		}()

		h.ServeHTTP(w, r)
	})
}

func (c *Components) unaryLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	t := time.Now()
	defer func() {
		d := time.Since(t)
		c.latency.Record(ctx, d.Milliseconds())

		c.Log.Trace().
			Str("method", info.FullMethod).
			Dur("dur", d).
			Msg("served")

	}()

	return handler(ctx, req)
}
