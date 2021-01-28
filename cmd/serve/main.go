package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	sitedata "go.seankhliao.com/com-seankhliao/v14"
	"go.seankhliao.com/com-seankhliao/v14/internal/serve"
	"k8s.io/klog/v2"
)

func main() {
	os.Exit(serve.Run(&Server{}))
}

type Server struct {
	f        fs.FS
	notfound http.Handler
}

func (s *Server) InitFlags(fs *flag.FlagSet) {}

func (s *Server) Setup(ctx context.Context, c *serve.Components) error {
	var err error
	s.f, err = fs.Sub(sitedata.S, "public")
	if err != nil {
		return fmt.Errorf("subfs public/: %w", err)
	}
	s.notfound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regexp.MustCompile(`/blog/\d{4}-\d{2}-\d{2}-.*/`).MatchString(r.URL.Path) {
			http.Redirect(w, r, "/blog/1"+r.URL.Path[6:], http.StatusMovedPermanently)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		f, err := s.f.Open("404.html")
		if err != nil {
			klog.ErrorS(err, "open 404.html")
		}
		_, err = io.Copy(w, f)
		if err != nil {
			klog.ErrorS(err, "write 404.html")
		}
	})
	fs.WalkDir(s.f, ".", func(op string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		p := op
		if strings.HasSuffix(p, "index.html") {
			p = p[:len(p)-10]
		} else if strings.HasSuffix(p, ".html") {
			p = p[:len(p)-5] + "/"
		}
		p = "/" + p
		c.Mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != p {
				klog.InfoS("not found", "path", r.URL.Path, "expected", p)
				s.notfound.ServeHTTP(w, r)
				return
			}

			setHeaders(w)
			var ct, cc string
			switch path.Ext(op) {
			case ".css":
				ct = "text/css"
				cc = "max-age=2592000"
			case ".js":
				ct = "application/javascript"
				cc = "max-age=2592000"
			case ".jpg", ".jpeg":
				ct = "image/jpeg"
				cc = "max-age=2592000"
			case ".png":
				ct = "image/png"
				cc = "max-age=2592000"
			case ".webp":
				ct = "image/webp"
				cc = "max-age=2592000"
			case ".svg":
				ct = "image/svg+xml"
				cc = "max-age=2592000"
			case ".json":
				ct = "application/json"
			case ".otf", ".ttf", ".woff", ".woff2":
				ct = "font/" + path.Ext(op)
			case ".html":
				ct = "text/html"
			}
			if ct != "" {
				w.Header().Set("content-type", ct)
			}
			if cc != "" {
				w.Header().Set("cache-control", cc)
			}

			f, err := s.f.Open(op)
			if err != nil {
				klog.ErrorS(err, "open", "file", op, "path", p)
			}
			_, err = io.Copy(w, f)
			if err != nil {
				klog.ErrorS(err, "write", "file", op, "path", p)
			}
		})
		return nil
	})

	klog.InfoS("setup complete")
	return nil
}
