package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.seankhliao.com/usvc"
)

const (
	name = "go.seankhliao.com/com-seankhliao"
)

func main() {
	os.Exit(usvc.Exec(context.Background(), &Server{}, os.Args))
}

type Server struct {
	dir      string
	notfound http.Handler

	log    zerolog.Logger
	tracer trace.Tracer

	page *prometheus.CounterVec
}

func (s *Server) Flags(fs *flag.FlagSet) {
	fs.StringVar(&s.dir, "dir", "public", "directory to serve")
}

func (s *Server) Setup(ctx context.Context, u *usvc.USVC) error {
	s.log = u.Logger
	s.tracer = global.Tracer(name)
	s.page = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "seankhliao_com_page_requests",
	}, []string{"page"})

	notfound, _ := ioutil.ReadFile(path.Join(s.dir, "404.html"))
	s.notfound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(notfound)
	})

	u.ServiceMux.Handle("/", s)
	return nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := s.tracer.Start(r.Context(), "handle-request")
	defer span.End()

	u, f := r.URL.Path, ""
	switch {
	case strings.HasSuffix(u, "/") && exists(path.Join(s.dir, u[:len(u)-1]+".html")):
		f = path.Join(s.dir, u[:len(u)-1]+".html")
	case strings.HasSuffix(u, "/") && exists(path.Join(s.dir, u, "index.html")):
		f = path.Join(s.dir, u, "index.html")
	case strings.HasSuffix(u, "/"):
		if s.notfound != nil {
			s.notfound.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		return
	case !strings.HasSuffix(u, ".html") && exists(path.Join(s.dir, u)):
		f = path.Join(s.dir, u)
	default:
		http.Redirect(w, r, canonical(u), http.StatusMovedPermanently)
		return
	}

	setHeaders(w)
	switch path.Ext(f) {
	case ".otf", ".ttf", ".woff", ".woff2", ".css", ".png", ".jpg", ".jpeg", ".webp", ".json", ".js":
		w.Header().Set("cache-control", `max-age=2592000`)
	}

	http.ServeFile(w, r, f)

	s.page.WithLabelValues(u).Add(1)
}
