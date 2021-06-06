package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-logr/logr"
	"go.seankhliao.com/w/v16/internal/render"
	"go.seankhliao.com/w/v16/internal/webserver"
	"k8s.io/klog/v2/klogr"
)

func main() {
	var o webserver.Options
	var dir, baseurl string
	flag.StringVar(&dir, "dir", "public", "path to directory to serve")
	flag.StringVar(&baseurl, "url", "https://arch.seankhliao.com", "base url for canonicalization")
	o.InitFlags(flag.CommandLine)
	flag.Parse()

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	l := klogr.New()

	m, err := newHttp(l, dir, baseurl)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	o.Logger = l
	o.Handler = m

	webserver.New(ctx, &o).Run(ctx)
}

func newHttp(l logr.Logger, dir, baseurl string) (*http.ServeMux, error) {
	tmpdir, err := os.MkdirTemp("", "singlepage")
	if err != nil {
		return nil, fmt.Errorf("temp dir: %w", err)
	}
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		if d.IsDir() {
			err = os.MkdirAll(filepath.Join(tmpdir, rel), 0o755)
			return err
		}
		if filepath.Ext(rel) == ".md" {
			rel = strings.TrimSuffix(rel, ".md") + ".html"
		}
		_, err = render.ProcessFile(p, filepath.Join(tmpdir, rel), baseurl, true, true)

		return err
	})

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(tmpdir)))
	return mux, nil
}
