package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.seankhliao.com/stream"
	"go.seankhliao.com/usvc"
	"google.golang.org/grpc"
)

const (
	name = "com-seankhliao"
)

func main() {
	var s Server

	usvc.Run(context.Background(), name, &s, false)
}

type Server struct {
	dir        string
	notfound   http.Handler
	streamAddr string
	client     stream.StreamClient
	cc         *grpc.ClientConn

	page metric.Int64Counter

	log zerolog.Logger
}

func (s *Server) Flag(fs *flag.FlagSet) {
	fs.StringVar(&s.dir, "dir", "public", "directory to serve")
	fs.StringVar(&s.streamAddr, "addr.stream", "stream:80", "url to connect to stream")
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

	var err error
	s.cc, err = grpc.Dial(s.streamAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("connect to stream: %w", err)
	}
	s.client = stream.NewStreamClient(s.cc)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.cc.Close()
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	remote := r.Header.Get("x-forwarded-for")
	if remote == "" {
		remote = r.RemoteAddr
	}

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

	httpRequest := &stream.HTTPRequest{
		Timestamp: time.Now().Format(time.RFC3339),
		Method:    r.Method,
		Domain:    r.Host,
		Path:      r.URL.Path,
		Remote:    remote,
		UserAgent: r.UserAgent(),
		Referrer:  r.Referer(),
	}

	_, err := s.client.LogHTTP(ctx, httpRequest)
	if err != nil {
		s.log.Error().Err(err).Msg("write to stream")
	}
}
