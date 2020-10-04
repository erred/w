package usvc

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.seankhliao.com/apis/saver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

const (
	name = "go.seankhliao.com/usvc"
)

type Service interface {
	Flags(fs *flag.FlagSet)
	Setup(ctx context.Context, usvc *USVC) error
}

type USVC struct {
	// external
	Logger zerolog.Logger
	// tracer is global
	// Tracer trace.Tracer
	// metrics is global

	ServiceMux    *http.ServeMux
	MetricMux     *http.ServeMux
	Live          HealthProbe
	Ready         HealthProbe
	ServiceServer *http.Server
	MetricServer  *http.Server
	GRPCServer    *grpc.Server

	// internal
	log    zerolog.Logger
	tracer trace.Tracer
}

// Exec runs the server
// returning an os.Exit code
//
// Globals: prometheus metrics, otel tracer
func Exec(ctx context.Context, svc Service, args []string) int {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		<-sigc
		cancel()
	}()

	var (
		serviceAddr string
		metricAddr  string
		loggerOpts  LoggerOpts
		metricOpts  MetricOpts
		tracerOpts  TracerOpts
		saverOpts   SaverOpts
		tlsOpts     TLSOpts
	)

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.StringVar(&serviceAddr, "addr", ":8080", "service listen address")
	fs.StringVar(&metricAddr, "addr.metric", ":8000", "metric listen address")
	loggerOpts.Flags(fs)
	metricOpts.Flags(fs)
	tracerOpts.Flags(fs)
	saverOpts.Flag(fs)
	tlsOpts.Flags(fs)
	svc.Flags(fs)
	fs.Parse(args[1:])

	usvc := USVC{
		Logger:     loggerOpts.Logger(true),
		ServiceMux: http.NewServeMux(),
		MetricMux:  http.NewServeMux(),
	}

	usvc.log = usvc.Logger.With().Str("module", "usvc").Logger()
	usvc.Live, usvc.Ready = metricOpts.Metrics(usvc.MetricMux)

	traceShut, err := tracerOpts.Tracer()
	if err != nil {
		usvc.log.Error().Err(err).Msg("init tracer")
		return 1
	}
	defer traceShut()
	usvc.tracer = global.Tracer(name)

	tlsConf, err := tlsOpts.Config()
	if err != nil {
		usvc.log.Error().Err(err).Msg("init tls")
		return 1
	}

	var do grpc.DialOption
	if tlsConf.RootCAs != nil {
		do = grpc.WithTransportCredentials(credentials.NewTLS(tlsConf))
	} else {
		do = grpc.WithInsecure()
	}
	saverClient, saverShut, err := saverOpts.Saver(grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor(usvc.tracer)), do)
	if err != nil {
		usvc.log.Error().Err(err).Msg("init saver")
		return 1
	}
	defer saverShut()

	// middleware
	latency := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "usvc_service_latency_s",
	}, []string{"proto"})

	var sos []grpc.ServerOption
	if len(tlsConf.Certificates) > 0 {
		sos = append(sos, grpc.Creds(credentials.NewTLS(tlsConf)))
	}
	sos = append(sos, grpc.ChainUnaryInterceptor(
		otelgrpc.UnaryServerInterceptor(global.Tracer(name)),
		unaryLogMiddleware(
			usvc.tracer,
			usvc.log,
			latency,
		),
	))
	usvc.GRPCServer = grpc.NewServer(sos...)

	handler := otelhttp.NewHandler(
		httpLogMiddleWare(
			usvc.tracer,
			usvc.log,
			latency,
			saverClient,
			corsAllowAll(usvc.ServiceMux),
		),
		"otelhttp",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)
	// handler := usvc.ServiceMux
	// _ = saverClient
	usvc.ServiceServer = &http.Server{
		Addr:              serviceAddr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
		TLSConfig:         tlsConf,
		ErrorLog:          log.New(usvc.Logger.With().Str("module", "http.service").Logger(), "", 0),
	}

	usvc.MetricServer = &http.Server{
		Addr:              metricAddr,
		Handler:           usvc.MetricMux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ErrorLog:          log.New(usvc.Logger.With().Str("module", "http.metric").Logger(), "", 0),
	}

	err = svc.Setup(ctx, &usvc)
	if err != nil {
		usvc.log.Error().Err(err).Msg("service setup")
		return 1
	}

	errs := usvc.run(ctx, cancel)
	if err != nil {
		usvc.log.Error().Errs("errs", errs).Msg("service run")
		return 1
	}
	usvc.log.Info().Msg("service shutdown cleanly")
	return 0
}

func (u *USVC) run(ctx context.Context, cancel func()) []error {
	errc := make(chan error)
	g := len(u.GRPCServer.GetServiceInfo()) != 0
	if g {
		go func() {
			u.log.Info().Str("addr", u.ServiceServer.Addr).Msg("starting service grpc endpoint")
			lis, err := net.Listen("tcp", u.ServiceServer.Addr)
			if err != nil {
				cancel()
				errc <- err
				return
			}
			err = u.GRPCServer.Serve(lis)
			cancel()
			errc <- err
		}()
	} else {
		go func() {
			u.log.Info().Str("addr", u.ServiceServer.Addr).Msg("starting service http endpoint")
			var err error
			if u.ServiceServer.TLSConfig != nil && len(u.ServiceServer.TLSConfig.Certificates) > 0 {
				err = u.ServiceServer.ListenAndServeTLS("", "")
			} else {
				err = u.ServiceServer.ListenAndServe()
			}
			cancel()
			errc <- err
		}()
	}

	go func() {
		u.log.Info().Str("addr", u.MetricServer.Addr).Msg("starting metrics endpoint")
		err := u.MetricServer.ListenAndServe()
		cancel()
		errc <- err
	}()

	<-ctx.Done()
	u.log.Info().Err(ctx.Err()).Msg("starting shutdown sequence")

	sdctx := context.Background()
	go func() {
		errc <- u.MetricServer.Shutdown(sdctx)
	}()
	if g {
		go func() {
			u.GRPCServer.GracefulStop()
			errc <- nil
		}()
	} else {
		go func() {
			errc <- u.ServiceServer.Shutdown(sdctx)
		}()
	}

	var errs []error
	for i := 0; i < 4; i++ {
		e := <-errc
		if errors.Is(e, http.ErrServerClosed) {
			// skip
		} else if e != nil {
			errs = append(errs, e)
		}
	}
	return errs
}

func httpLogMiddleWare(tracer trace.Tracer, log zerolog.Logger, latency *prometheus.HistogramVec, saverClient saver.SaverClient, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "httpLogMiddleWare")
		defer span.End()

		t := time.Now()
		defer func() {
			d := time.Since(t)
			latency.WithLabelValues(r.Proto).Observe(d.Seconds())

			remote := r.Header.Get("x-forwarded-for")
			if remote == "" {
				remote = r.RemoteAddr
			}
			log.Trace().
				Str("src", remote).
				Str("url", r.URL.String()).
				Str("user-agent", r.UserAgent()).
				Dur("dur", d).
				Msg("served")

			ctx, span := tracer.Start(ctx, "saverClient")
			defer span.End()
			_, err := saverClient.HTTP(ctx, &saver.HTTPRequest{
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
				log.Error().Err(err).Msg("send to saver")
			}
		}()

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func unaryLogMiddleware(tracer trace.Tracer, log zerolog.Logger, latency *prometheus.HistogramVec) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, span := tracer.Start(ctx, "unaryLogMiddleware")
		defer span.End()

		t := time.Now()
		defer func() {
			d := time.Since(t)
			latency.WithLabelValues("grpc").Observe(d.Seconds())

			var pa string
			p, ok := peer.FromContext(ctx)
			if ok {
				pa = p.Addr.String()
			}

			log.Trace().
				Str("src", pa).
				Str("method", info.FullMethod).
				Dur("dur", d).
				Msg("served")
		}()

		return handler(ctx, req)
	}
}

func corsAllowAll(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		case http.MethodGet, http.MethodPost:
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Max-Age", "86400")
			h.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}
