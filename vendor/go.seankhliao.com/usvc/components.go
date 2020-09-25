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
	"go.seankhliao.com/apis/saver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

type Components struct {
	Log    zerolog.Logger
	Meter  metric.Meter
	Tracer trace.Tracer

	// added outside
	HTTP *http.ServeMux
	GRPC *grpc.Server

	GRPCDialOptions []grpc.DialOption

	// internal
	latency     metric.Int64ValueRecorder
	flush       func()
	prom        *prometheus.Exporter
	saverConn   *grpc.ClientConn
	saverClient saver.SaverClient
}

func NewComponents(name, logLevel, logFormat, caCrt, saverURL string) (*Components, error) {
	c := &Components{}
	// Note: always return c, the so we have a logger

	c.logger(logLevel, logFormat)

	err := c.client(caCrt)
	if err != nil {
		return c, fmt.Errorf("setup client crt: %w", err)
	}

	err = c.tracer(name)
	if err != nil {
		return c, fmt.Errorf("install tracer: %w", err)
	}

	err = c.metric(name)
	if err != nil {
		return c, fmt.Errorf("install metric: %w", err)
	}

	err = c.saver(saverURL)
	if err != nil {
		return c, fmt.Errorf("initialize saver: %w", err)
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

func (c *Components) client(caCrt string) error {
	if caCrt != "" {
		creds, err := credentials.NewClientTLSFromFile(caCrt, "")
		if err != nil {
			return err
		}
		c.GRPCDialOptions = append(c.GRPCDialOptions, grpc.WithTransportCredentials(creds))
	}
	return nil
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
		"usvc_request_serve_latency_ms",
		metric.WithDescription("time to handle request"),
		metric.WithUnit(unit.Milliseconds),
	)
	return nil
}

func (c *Components) saver(saverURL string) error {
	if saverURL != "" && len(c.GRPCDialOptions) > 0 {
		var err error
		c.saverConn, err = grpc.Dial(saverURL, c.GRPCDialOptions...)
		if err != nil {
			return fmt.Errorf("saver dial %s: %w", saverURL, err)
		}
		c.saverClient = saver.NewSaverClient(c.saverConn)
	}
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

			if c.saverClient != nil {
				_, err := c.saverClient.HTTP(context.Background(), &saver.HTTPRequest{
					HttpRemote: &saver.HTTPRemote{
						Timestamp: t.Format(time.RFC3339),
						Remote:    remote,
						UserAgent: r.UserAgent(),
						Referrer:  r.Referer(),
					},
					Method: r.Method,
					Domain: r.Host,
					Path:   r.URL.Path,
				})
				if err != nil {
					c.Log.Error().Err(err).Msg("send to saver")
				}
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func (c *Components) unaryLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	t := time.Now()
	defer func() {
		d := time.Since(t)
		c.latency.Record(ctx, d.Milliseconds())

		var pa string
		p, ok := peer.FromContext(ctx)
		if ok {
			pa = p.Addr.String()
		}

		c.Log.Trace().
			Str("src", pa).
			Str("method", info.FullMethod).
			Dur("dur", d).
			Msg("served")
	}()

	return handler(ctx, req)
}
