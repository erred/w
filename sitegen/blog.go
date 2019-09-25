package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/russross/blackfriday/v2"
)

type BlogData struct {
	URL     string
	Date    string
	Title   string
	Desc    string
	Content string

	AbsURL string // for blog index
	// only for index page
	Posts []BlogData
}

// BlogOptions holds config needed for parsing blog posts
// expects a flat directory of files named: yyyy-mm-dd-some-titile.md
// evenything before the first "\n---\n" is used as a description
//
// dependencies on:
//      template named "blogpost"
//      template names "blogindex"
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

	var datas []BlogData

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
		fn := strings.TrimSuffix(base, ".md")

		b, err := ioutil.ReadFile(filepath.Join(o.Src, fi.Name()))
		if err != nil {
			log.Printf("BlogOptions.Exec read file %q: %q\n", fi.Name(), err)
			continue
		}
		bb := bytes.SplitN(b, []byte("\n---\n"), 2)
		var desc []byte
		if len(bb) == 2 {
			desc = bytes.TrimSpace(bb[0])
			b = bytes.TrimSpace(bb[1])
		}

		data := BlogData{
			URL:     filepath.Join(opt.host, o.Src, fn),
			AbsURL:  filepath.Join("/", o.Src, fn),
			Date:    strings.Join(parts[:3], "-"),
			Title:   strings.Join(parts[3:], " "),
			Desc:    string(desc),
			Content: string(blackfriday.Run(b, blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{HeadingLevelOffset: 3, Flags: blackfriday.CommonHTMLFlags})))),
		}
		datas = append(datas, data)

		var buf bytes.Buffer
		err = opt.T.ExecuteTemplate(&buf, "blogpost", data)
		if err != nil {
			log.Printf("BlogOptions.Exec exec template %q: %w\n", fn, err)
			continue
		}

		dfn := filepath.Join(o.Dst, o.Src, fn+".html")
		os.MkdirAll(filepath.Dir(dfn), 0755)
		err = ioutil.WriteFile(dfn, buf.Bytes(), 0644)
		if err != nil {
			log.Printf("BlogOptions.Exec write %q: %w\n", dfn, err)
		}
	}

	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Date > datas[j].Date
	})

	data := BlogData{
		URL:   filepath.Join(opt.host, o.Src),
		Title: o.Src,
		Posts: datas,
	}

	var buf bytes.Buffer
	err = opt.T.ExecuteTemplate(&buf, "blogindex", data)
	if err != nil {
		return fmt.Errorf("BlogOptions.Exec exec template blogindex: %w", err)
	}

	dfn := filepath.Join(o.Dst, o.Src, "index.html")
	os.MkdirAll(filepath.Dir(dfn), 0755)
	err = ioutil.WriteFile(dfn, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("BlogOptions.Exec write index %q: %w", dfn, err)
	}

	return nil
}
