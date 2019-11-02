package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SitemapOptions holds config for generating a txt sitemap
type SitemapOptions struct {
	Src string
	Dst string
}

func NewSitemapOptions(args []string) *SitemapOptions {
	var o SitemapOptions
	f := flag.NewFlagSet("mod", flag.ExitOnError)
	f.StringVar(&o.Src, "src", "dst", "input directory")
	f.StringVar(&o.Dst, "dst", "dst/sitemap.txt", "output file")
	f.Parse(args)
	return &o
}

func (o *SitemapOptions) Exec(opt *Options) error {
	var files []string
	filepath.Walk(o.Src, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}
		if filepath.Ext(path) != ".html" {
			return nil
		}
		subpath, _ := filepath.Rel(o.Src, path)
		files = append(files, "https://"+filepath.Join(opt.host, canonicalURL(subpath)))
		return nil
	})

	sort.Strings(files)
	err := ioutil.WriteFile(o.Dst, []byte(strings.Join(files, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("SitemapOptions.Exec write %q: %w", o.Dst, err)
	}
	return nil
}
