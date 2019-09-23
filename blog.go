package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

type BlogData struct {
	URL     string
	Title   string
	Content string
}

// BlogOptions holds config needed for parsing blog posts
// expects a flat directory of files named: yyyy-mm-dd-some-titile.md
// dependencies on:
//      template named "blogpost"
type BlogOptions struct {
	Src string
	Dst string
}

func NewBlogOptions(args []string) *BlogOptions {
	var o BlogOptions
	f := flag.NewFlagSet("blog", flag.ExitOnError)
	f.StringVar(&o.Src, "src", "blog", "source directory")
	f.StringVar(&o.Dst, "dst", "dst", "output directory")
	f.Parse(args)
	return &o
}

func (o *BlogOptions) Exec(opt *Options) error {
	fis, err := ioutil.ReadDir(o.Src)
	if err != nil {
		return fmt.Errorf("BlogOptions.Exec read dir %q: %w", o.Src, err)
	}
	for _, fi := range fis {
		if fi.IsDir() || filepath.Ext(fi.Name()) != ".md" {
			continue
		}
		base := filepath.Base(fi.Name())
		parts := strings.Split(strings.TrimSuffix(base, ".md"), "-")
		if len(parts) < 4 {
			log.Printf("BlogOptions.Exec parse name %q, expected at least 4 parts, got %d\n", base, len(parts))
			continue
		}
		fn := base + ".html"

		b, err := ioutil.ReadFile(fi.Name())
		if err != nil {
			log.Printf("BlogOptions.Exec read file %q: %q\n", fi.Name(), err)
			continue
		}
		data := BlogData{
			URL:     filepath.Join(opt.host, "blog", fn),
			Title:   strings.Join(parts[3:], " "),
			Content: string(blackfriday.Run(b, blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{HeadingLevelOffset: 3, Flags: blackfriday.CommonHTMLFlags})))),
		}

		var buf bytes.Buffer
		err = opt.T.ExecuteTemplate(&buf, "blogpage", data)
		if err != nil {
			log.Printf("BlogOptions.Exec exec template %q: %w\n", fn, err)
			continue
		}

		dfn := filepath.Join(o.Dst, "blog", fn)
		err = ioutil.WriteFile(dfn, buf.Bytes(), 0644)
		if err != nil {
			log.Printf("BlogOptions.Exec write %q: %w\n", dfn, err)
		}
	}
	return fmt.Errorf("ErrNotImplemented")
}
