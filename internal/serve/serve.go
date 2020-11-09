package serve

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/klog/v2"
)

type Service interface {
	// must handle nil
	InitFlags(fs *flag.FlagSet)
	Setup(ctx context.Context, c *Components) error
}

type Components struct {
	Mux    *http.ServeMux
	Server *http.Server
}

func Run(svc Service) int {
	var addr string
	flag.StringVar(&addr, "addr", os.Getenv("PORT"), "address")
	klog.InitFlags(nil)
	svc.InitFlags(nil)
	flag.Parse()
	if addr == "" {
		addr = ":8080"
	} else if addr[0] != ':' {
		addr = ":" + addr
	}

	mux := http.NewServeMux()
	c := &Components{
		Mux: mux,
		Server: &http.Server{
			Addr:              addr,
			Handler:           corsAllowAll(mux),
			ReadHeaderTimeout: 10 * time.Second,
			MaxHeaderBytes:    1 << 20,
			// ErrorLog
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		<-sigc
		cancel()
	}()

	err := svc.Setup(ctx, c)
	if err != nil {
		klog.ErrorS(err, "service setup")
		return 2
	}

	errc := make(chan error)
	go func() {
		defer cancel()
		klog.InfoS("starting server", "addr", c.Server.Addr)
		err := c.Server.ListenAndServe()
		switch {
		case errors.Is(err, http.ErrServerClosed):
			close(errc)
		case err != nil:
			errc <- err
		default:
			close(errc)
		}
	}()

	<-ctx.Done()
	err = c.Server.Shutdown(context.Background())
	if err != nil {
		klog.ErrorS(err, "server shutdown")
		return 1
	}
	return 0
}

func corsAllowAll(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("server", "internal/serve/v13")
		w.Header().Set("hire-me", "http-header-hire@seankhliao.com")
		w.Header().Set("easter-egg", "🐇*(🍆-🪴)=🐇🥚")

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
