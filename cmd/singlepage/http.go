package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
	"go.seankhliao.com/w/v15/internal/render"
)

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
