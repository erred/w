package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"go.seankhliao.com/com-seankhliao/v13/internal/serve"
	"k8s.io/klog/v2"
)

func main() {
	os.Exit(serve.Run(&Server{}))
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

func (s *Server) Setup(ctx context.Context, c *serve.Components) error {
	notfound, _ := ioutil.ReadFile(path.Join(s.dir, "404.html"))
	s.notfound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regexp.MustCompile(`/blog/\d{4}-\d{2}-\d{2}-.*/`).MatchString(r.URL.Path) {
			http.Redirect(w, r, "/blog/1"+r.URL.Path[6:], http.StatusMovedPermanently)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(notfound)
	})

	c.Mux.Handle("/", s)
	klog.InfoS("setup complete", "dir", s.dir, "notfound", notfound != nil)
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
