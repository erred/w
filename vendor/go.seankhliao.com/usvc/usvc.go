package usvc

import (
	"context"
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Service only mandates a name,
// everythign else is optional
type Service interface {
	// Flagger is for registering flags
	Flag(fs *flag.FlagSet)
	// Register is the time to
	// get the logger, tracer, metric;
	// add routes, services;
	// start background services
	Register(c *Components) error
	// Shutdown an optional interface to implement
	// to be called in the graceful shutdown sequence
	Shutdown(ctx context.Context) error
}

func Run(ctx context.Context, name string, server Service, grpcsvc bool) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	metricAddr := fs.String("addr.metric", ":8000", "metric address")
	serviceAddr := fs.String("addr", ":8080", "service address")
	tlsCert := fs.String("tls.crt", "", "tls cert file")
	tlsKey := fs.String("tls.key", "", "tls key file")
	logLevel := fs.String("log.lvl", "trace", "log level: trace, debug, info, error")
	logFormat := fs.String("log.fmt", "json", "log format: logfmt, json")
	server.Flag(fs)
	fs.Parse(os.Args[1:])

	c, err := NewComponents(name, *logLevel, *logFormat)
	if err != nil {
		c.Log.Error().Err(err).Msg("setup components")
		return
	}

	tlsConf, gopts, err := tlsSetup(*tlsCert, *tlsKey)
	if err != nil {
		c.Log.Error().Err(err).Str("crt", *tlsCert).Str("key", *tlsKey).Msg("setup tls")
		return
	}

	// metrics endpoint
	mmux := http.NewServeMux() // metrics http
	mmux.Handle("/metrics", c.prom)
	mmux.HandleFunc("/health", okhandler)
	pprofHandlers(mmux)
	msrv := httpServer(*metricAddr, mmux, nil)

	// grpc server endpoint
	gopts = append(gopts, grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor(c.Tracer), c.unaryLog))
	c.GRPC = grpc.NewServer(gopts...)

	// http server endpoint
	c.HTTP = http.NewServeMux() // service http
	hhandler := otelhttp.NewHandler(c.httpLog(corsAllowAll(c.HTTP)), name, otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))
	hsrv := httpServer(*serviceAddr, hhandler, tlsConf)

	// register other services
	err = server.Register(c)
	if err != nil {
		c.Log.Error().Err(err).Msg("register service")
		return
	}

	// run
	svc := "http"
	if grpcsvc {
		svc = "grpc"
	}
	c.Log.Info().Str(svc, *serviceAddr).Bool("tls", tlsConf != nil).Str("metric", *metricAddr).Msg("starting server")
	errc, done := startServers(ctx, msrv, hsrv, c.GRPC, tlsConf, grpcsvc)
	<-done
	shutdownServers(errc, server, msrv, hsrv, c.GRPC, grpcsvc)

	var errs []error
	for err := range errc {
		if err != nil {
			errs = append(errs, err)
		}
	}
	c.Log.Error().Errs("shutdown", errs).Msg("exit")
}

func tlsSetup(tlsCert, tlsKey string) (*tls.Config, []grpc.ServerOption, error) {
	if tlsCert != "" && tlsKey != "" {
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			return nil, nil, err
		}
		tlsConf := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
		}
		gopts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(tlsConf))}
		return tlsConf, gopts, nil
	}
	return nil, nil, nil
}

func okhandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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

func startServers(ctx context.Context, msrv, hsrv *http.Server, gsrv *grpc.Server, tlsConf *tls.Config, grpcsvc bool) (chan error, <-chan struct{}) {
	errc := make(chan error)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		<-sigc
		cancel()
	}()
	go func() {
		err := msrv.ListenAndServe()
		cancel()
		errc <- err
	}()
	if grpcsvc {
		go func() {
			var err error
			lis, err := net.Listen("tcp", "")
			if err != nil {
				cancel()
				errc <- err
				return
			}
			err = gsrv.Serve(lis)
			cancel()
			errc <- err
		}()
	} else {
		go func() {
			var err error
			if tlsConf == nil {
				err = hsrv.ListenAndServe()
			} else {
				err = hsrv.ListenAndServeTLS("", "")
			}
			cancel()
			errc <- err
		}()
	}
	return errc, ctx.Done()
}

func shutdownServers(errc chan error, server Service, msrv, hsrv *http.Server, gsrv *grpc.Server, grpcsvc bool) {
	sdctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		errc <- server.Shutdown(sdctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		errc <- msrv.Shutdown(sdctx)
	}()
	if grpcsvc {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gsrv.GracefulStop()
		}()
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errc <- hsrv.Shutdown(sdctx)
		}()
	}

	go func() {
		wg.Wait()
		close(errc)
	}()
}

func pprofHandlers(mmux *http.ServeMux) {
	mmux.HandleFunc("/debug/pprof/", pprof.Index)
	mmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mmux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mmux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mmux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mmux.Handle("/debug/pprof/block", pprof.Handler("block"))
}

func httpServer(addr string, h http.Handler, tlsConf *tls.Config) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: h,
		// ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		// WriteTimeout:      10 * time.Second,
		// IdleTimeout:       60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConf,
	}
}
