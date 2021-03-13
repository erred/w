package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.seankhliao.com/w/v15/internal/static"
	"k8s.io/klog/v2/klogr"
)

func main() {
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	l := klogr.New()
	otel.SetErrorHandler(otelErrorHandler{l})

	shutdown, mh, err := o11y(ctx)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}
	defer shutdown(context.Background()) // separate shutdown context

	NewW(l, mh).run(ctx)
}

type W struct {
	l logr.Logger

	m, h *http.Server
}

func NewW(l logr.Logger, mh http.Handler) *W {
	// metrics server
	mmux := http.NewServeMux()
	mmux.HandleFunc("/debug/pprof/", pprof.Index)
	mmux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mmux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mmux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mmux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mmux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mmux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mmux.HandleFunc("/healthyz", func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte("live")) })
	mmux.HandleFunc("/readyz", func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte("ready")) })
	mmux.Handle("/metrics", mh)

	// http server
	hmux, err := newHttp(l, static.S)
	if err != nil {
		l.Error(err, "setup http")
		os.Exit(1)
	}

	w := &W{
		l: l,
		h: &http.Server{
			Addr:              ":8080",
			Handler:           otelhttp.NewHandler(corsAllowAll(hmux), "serve"),
			ReadHeaderTimeout: 10 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
		m: &http.Server{
			Addr:              ":8090",
			Handler:           mmux,
			ReadHeaderTimeout: 10 * time.Second,
			MaxHeaderBytes:    1 << 10,
		},
	}
	return w
}

func (w W) run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, s := range []*http.Server{w.m, w.h} {
		wg.Add(2)

		// handle shutdown
		go func(s *http.Server) {
			defer wg.Done()
			<-ctx.Done()
			s.Shutdown(context.Background())
		}(s)

		// start servers
		go func(s *http.Server) {
			defer wg.Done()
			defer cancel()
			w.l.Info("starting", "addr", s.Addr)
			err := s.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				w.l.Error(err, "server closed")
			}
		}(s)
	}
}

type otelErrorHandler struct {
	l logr.Logger
}

func (o otelErrorHandler) Handle(err error) {
	o.l.Error(err, "otel error")
}
