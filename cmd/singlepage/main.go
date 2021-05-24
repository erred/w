package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2/klogr"
)

func main() {
	var dir, baseurl string
	flag.StringVar(&dir, "dir", "public", "path to directory to serve")
	flag.StringVar(&baseurl, "url", "https://arch.seankhliao.com", "base url for canonicalization")
	flag.Parse()

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	l := klogr.New()

	New(l, dir, baseurl).Run(ctx)
}

type SP struct {
	l logr.Logger

	m, h *http.Server
}

func New(l logr.Logger, dir, baseurl string) *SP {
	hmux, err := newHttp(l, dir, baseurl)
	if err != nil {
		l.Error(err, "setup http")
		os.Exit(1)
	}

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
	mmux.HandleFunc("/healthz", func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte("live")) })
	mmux.HandleFunc("/readyz", func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte("ready")) })

	return &SP{
		l: l,
		h: &http.Server{
			Addr:              ":8080",
			Handler:           hmux,
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
}

func (sp *SP) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, s := range []*http.Server{sp.m, sp.h} {
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
			sp.l.Info("starting", "addr", s.Addr)
			err := s.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				sp.l.Error(err, "server closed")
			}
		}(s)
	}
}
