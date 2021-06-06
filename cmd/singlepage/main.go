package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.seankhliao.com/w/v16/render"
	"go.seankhliao.com/w/v16/webserver"
	"k8s.io/klog/v2/klogr"
)

func main() {
	var wo webserver.Options
	var ro render.Options
	var fn string
	flag.StringVar(&fn, "file", "index.md", "file to serve")
	flag.StringVar(&ro.Data.GTMID, "gtm", "", "Google Tag Manager ID for analytics")
	flag.StringVar(&ro.Data.URLCanonical, "canonical", "https://arch.seankhliao.com", "canonical base url")
	flag.BoolVar(&ro.Data.Compact, "compact", true, "compact header")
	flag.BoolVar(&ro.MarkdownSkip, "raw", false, "skip markdown processing")
	wo.InitFlags(flag.CommandLine)
	flag.Parse()

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	l := klogr.New()

	m, err := newHttp(&ro, fn)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	wo.Logger = l
	wo.Handler = m

	webserver.New(ctx, &wo).Run(ctx)
}

func newHttp(ro *render.Options, fn string) (*http.ServeMux, error) {
	fin, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", fn, err)
	}

	buf := &bytes.Buffer{}
	err = render.Render(ro, buf, fin)
	if err != nil {
		return nil, fmt.Errorf("render: %w", err)
	}
	b := buf.Bytes()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	})
	return mux, nil
}
