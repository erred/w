package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
	"golang.org/x/tools/blog/atom"
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

		dfn := filepath.Join(o.Dst, o.Src, fn+".html")
		err = writeTemplate(opt.T, "blogpost", dfn, data)
		if err != nil {
			log.Printf("BlogOptions.Exec write blogpost %q: %w\n", dfn, err)
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

	dfn := filepath.Join(o.Dst, o.Src, "index.html")
	err = writeTemplate(opt.T, "blogindex", dfn, data)
	if err != nil {
		return fmt.Errorf("BlogOptions.Exec write blogindex: %w", err)
	}

	me := &atom.Person{
		Name:  "Sean Liao",
		URI:   "https://seankhliao.com",
		Email: "blog-atom@seankhliao.com",
	}

	fd := atom.Feed{
		Title: "seankhliao's stream of consciousness",
		ID:    "tag:seankhliao.com,2019:seankhliao.com",
		Link: []atom.Link{
			{
				Rel:  "self",
				Href: "https://seankhliao.com/feed.atom",
				Type: "application/atom+xml",
			}, {
				Rel:  "alternate",
				Href: "https://seankhliao.com/blog",
				Type: "text/html",
			},
		},
		Updated: atom.Time(time.Now()),
		Author:  me,
	}
	for _, bp := range datas {
		fd.Entry = append(fd.Entry, &atom.Entry{
			Title: bp.Title,
			Link: []atom.Link{
				{
					Rel:  "alternate",
					Href: bp.URL,
					Type: "text/html",
				},
			},
			ID:        "https://" + bp.URL,
			Published: atom.TimeStr(bp.Date + "T00:00:00Z"),
			Updated:   atom.TimeStr(bp.Date + "T00:00:00Z"),
			Author:    me,
			Summary: &atom.Text{
				Type: "text",
				Body: bp.Desc,
			},
		})
	}

	f, err := os.Create(filepath.Join(o.Dst, "feed.atom"))
	if err != nil {
		return fmt.Errorf("BlogOptions.Exec open feed.atom: %w", err)
	}
	defer f.Close()
	e := xml.NewEncoder(f)
	e.Indent("", "    ")
	err = e.Encode(fd)
	if err != nil {
		return fmt.Errorf("BlogOptions.Exec marshal atom: %w", err)
	}

	return nil
}
