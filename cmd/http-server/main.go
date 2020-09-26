package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.seankhliao.com/usvc"
)

const (
	name = "com-seankhliao"
)

func main() {
	var s Server

	usvc.Run(context.Background(), name, &s, false)
}

type Server struct {
	dir      string
	notfound http.Handler

	page metric.Int64Counter

	log zerolog.Logger
}

func (s *Server) Flag(fs *flag.FlagSet) {
	fs.StringVar(&s.dir, "dir", "public", "directory to serve")
}

func (s *Server) Register(c *usvc.Components) error {
	s.log = c.Log
	s.page = metric.Must(global.Meter(os.Args[0])).NewInt64Counter(
		"page_hit",
		metric.WithDescription("hits per page"),
	)

	notfound, _ := ioutil.ReadFile(path.Join(s.dir, "404.html"))
	s.notfound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(notfound)
	})

	c.HTTP.Handle("/", s)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
}
