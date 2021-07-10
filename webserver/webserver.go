package webserver

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.seankhliao.com/w/v16/internal/stdlog"
)

type Options struct {
	AdmAddr      string
	AppAddr      string
	OtlpEndpoint string
	Logger       logr.Logger
	Handler      http.Handler
}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	o.InitFlags(fs)
	return &o
}

func (o *Options) InitFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.AdmAddr, "adm.addr", ":8090", "listen address for admin")
	fs.StringVar(&o.AppAddr, "web.addr", ":8080", "listen address for main app")
	fs.StringVar(&o.OtlpEndpoint, "otlp.endpoint", "", "otlp grpc endpoint")
}

type Server struct {
	log logr.Logger

	adm *http.Server
	app *http.Server

	// shutdown
	sd *shutdown
}

func New(ctx context.Context, o *Options) *Server {
	adm := http.NewServeMux()

	// pprof
	adm.HandleFunc("/debug/pprof/", pprof.Index)
	adm.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	adm.Handle("/debug/pprof/block", pprof.Handler("block"))
	adm.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	adm.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	adm.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	adm.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	adm.HandleFunc("/debug/pprof/profile", pprof.Profile)
	adm.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	adm.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	adm.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// healthchecks
	adm.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	adm.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	// o11y
	otel.SetErrorHandler(&errhandler{o.Logger.WithName("otel")})
	sd, err := o11y(ctx, o.OtlpEndpoint)
	if err != nil {
		o.Logger.Error(err, "setup o11y")
		os.Exit(1)
	}

	// adm.Handle("/metrics", mh)

	return &Server{
		log: o.Logger.WithName("webserver"),
		adm: &http.Server{
			Addr:              o.AdmAddr,
			Handler:           adm,
			ReadHeaderTimeout: 10 * time.Second,
			MaxHeaderBytes:    1 << 10,
			ErrorLog:          stdlog.New(o.Logger.WithName("admsvr"), false),
		},
		app: &http.Server{
			Addr:              o.AppAddr,
			Handler:           otelhttp.NewHandler(o.Handler, "appsvr"),
			ReadHeaderTimeout: 10 * time.Second,
			MaxHeaderBytes:    1 << 20,
			BaseContext: func(net.Listener) context.Context {
				ctx := context.Background()
				ctx = logr.NewContext(ctx, o.Logger.WithName("handler"))
				return ctx
			},
			ErrorLog: stdlog.New(o.Logger.WithName("appsvr"), false),
		},
		sd: sd,
	}
}

func (s *Server) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(3)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		<-ctx.Done()

		s.adm.Shutdown(context.Background())
		s.app.Shutdown(context.Background())
		if s.sd != nil {
			s.sd.Shutdown(context.Background())
		}
	}()

	go func() {
		defer wg.Done()
		defer cancel()

		s.log.Info("starting", "svr", "adm", "addr", s.adm.Addr)
		err := s.adm.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error(err, "shutdown", "svr", "adm")
		}
	}()

	go func() {
		defer wg.Done()
		defer cancel()

		s.log.Info("starting", "svr", "app", "addr", s.app.Addr)
		err := s.app.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error(err, "shutdown", "svr", "app")
		}
	}()
}

type errhandler struct {
	l logr.Logger
}

func (e *errhandler) Handle(err error) {
	if err != nil {
		e.l.Error(err, "")
	}
}
