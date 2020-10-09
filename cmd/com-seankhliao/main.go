package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	name = "go.seankhliao.com/com-seankhliao"
)

func main() {
	os.Exit(serve(&Server{}))
}

type Server struct {
	dir      string
	notfound http.Handler
}

func (s *Server) InitFlags(fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}
	fs.StringVar(&s.dir, "dir", "public", "directory to serve")
}

func (s *Server) Setup(ctx context.Context, c *Components) error {
	notfound, _ := ioutil.ReadFile(path.Join(s.dir, "404.html"))
	s.notfound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(notfound)
	})

	c.Mux.Handle("/", s)
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
