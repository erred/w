package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
)

type Http struct {
	l        logr.Logger
	f        fs.FS
	notfound http.Handler
}

type HttpMux interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
}

func newHttp(l logr.Logger, sitedata fs.FS) (*http.ServeMux, error) {
	h := &Http{
		l: l,
		f: sitedata,
	}
	h.notfound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if regexp.MustCompile(`/blog/\d{4}-\d{2}-\d{2}-.*/`).MatchString(r.URL.Path) {
			http.Redirect(w, r, "/blog/1"+r.URL.Path[6:], http.StatusMovedPermanently)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		h.serveFile(w, "404.html")
	})

	mux := http.NewServeMux()
	err := fs.WalkDir(h.f, ".", func(op string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		// canonical path
		p := canonicalPath(op)

		l.Info("registering", "path", p)
		mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			// handle unknown paths
			if r.URL.Path != p {
				l.Error(errors.New("path mismatch"), "not found", "expected", p, "got", r.URL.Path)
				h.notfound.ServeHTTP(w, r)
				return
			}

			setHeaders(w, op)
			h.serveFile(w, op)
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("newHttp: walk fs: %w", err)
	}
	return mux, nil
}

func (h Http) serveFile(w http.ResponseWriter, p string) {
	file, err := h.f.Open(p)
	if err != nil {
		h.l.Error(err, "open", "file", p)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		h.l.Error(err, "copy", "file", p)
	}
}

func canonicalPath(p string) string {
	if strings.HasSuffix(p, "index.html") {
		p = p[:len(p)-10]
	} else if strings.HasSuffix(p, ".html") {
		p = p[:len(p)-5] + "/"
	}
	return "/" + p
}

func setHeaders(w http.ResponseWriter, op string) {
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
}

func corsAllowAll(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("svr", "go.seankhliao.com/w/v14")
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
