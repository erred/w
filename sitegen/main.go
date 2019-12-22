package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
	"text/template"
)

var (
	defaultSrc       = "src"
	defaultDst       = "dst"
	defaultTemplates = "templates/*"
	defaultFonts     = "https://fonts.googleapis.com/css?family=Inconsolata:400,700|Lora:400,700&display=swap"
	defaultAMPPrefux = "amp"
	defaultBaseURL   = "https://seankhliao.com"
	defaultGAID      = "UA-114337586-1"

	defaultImgArgs = [][]string{
		[]string{
			"-background", "none",
			"-density", "1200",
			"-resize", "1920x1080",
			"imgs/map.svg",
			"-write", "dst/map.png", "dst/map.webp", "dst/map.jpg",
		},
		[]string{
			"-flatten",
			"(", "+clone", "-resize", "512x512", "-quality", "60", "-write", "dst/icon-512.png", "+delete", ")",
			"(", "+clone", "-resize", "192x192", "-quality", "60", "-write", "dst/icon-192.png", "+delete", ")",
			"(", "+clone", "-resize", "128x128", "-quality", "60", "-write", "dst/icon-128.png", "+delete", ")",
			"(", "+clone", "-resize", "64x64", "-quality", "60", "-write", "dst/icon-64.png", "+delete", ")",
			"(", "+clone", "-resize", "48x48", "-quality", "60", "-write", "dst/icon-48.png", "+delete", ")",
			"(", "+clone", "-resize", "32x32", "-quality", "60", "-write", "dst/icon-32.png", "+delete", ")",
			"(", "+clone", "-resize", "16x16", "-quality", "60", "-write", "dst/icon-16.png", "+delete", ")",
			"-resize", "32x32", "-write", "dst/favicon.ico",
		},
	}
)

func main() {
	o := newOptions()

	// err := o.parseTemplates()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	//
	// err = o.getFonts()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	//
	// err = o.convertImgs()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }

	o.processPages()
}

type options struct {
	src          string
	dst          string
	templateGlob string
	fontURL      string
	ampPrefix    string
	baseURL      string
	gaID         string

	templates *template.Template
}

func newOptions() *options {
	o := &options{}
	flag.StringVar(&o.src, "src", defaultSrc, "source directory")
	flag.StringVar(&o.dst, "dst", defaultDst, "output directory")
	flag.StringVar(&o.templateGlob, defaultTemplates, "templates", "template glob")
	flag.StringVar(&o.fontURL, "fonts", defaultFonts, "fonts to retreive")
	flag.StringVar(&o.ampPrefix, "amp", defaultAMPPrefux, "amp path prefix")
	flag.StringVar(&o.baseURL, "host", defaultBaseURL, "url base")
	flag.StringVar(&o.gaID, "ga", defaultGAID, "google analytics ID")
	flag.Parse()
	return o
}

func (o *options) parseTemplates() error {
	var err error
	o.templates, err = template.ParseGlob(o.templateGlob)
	if err != nil {
		return fmt.Errorf("options.parseTemplates: %w", err)
	}
	return nil
}

func (o *options) getFonts() error {
	res, err := http.Get(o.fontURL)
	if err != nil {
		return fmt.Errorf("options.getFonts: %w", err)
	} else if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("options.getFonts: %d %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()
	buf := bytes.NewBufferString(`{{ define "fontcss" }}`)
	buf.ReadFrom(res.Body)
	buf.WriteString(`{{ end }}`)

	o.templates, err = o.templates.New("fontcss").Parse(buf.String())
	if err != nil {
		return fmt.Errorf("options.getFonts: %w", err)
	}
	return nil
}

func (o *options) convertImgs() error {
	for i, imgArgs := range defaultImgArgs {
		out, err := exec.Command("convert", imgArgs...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("options.convertImgs: %d: %w\n%s", i, err, out)
		}
	}
	return nil
}

func (o *options) processPages() error {
	fis, err := ioutil.ReadDir(o.src)
	if err != nil {
		return fmt.Errorf("options.parsePages: %w", err)
	}

	var wg sync.WaitGroup
	collect := make(chan page)
	for _, fi := range fis {
		wg.Add(1)
		if filepath.Ext(fi.Name()) == ".md" {
			go func(fn string) {
				defer wg.Done()
				p, err := o.processPage(fn)
				if err != nil {
					log.Println("options.parsePages: ", err)
				}
				collect <- p
			}(fi.Name())
		} else {
			// copy file
		}
	}
	go func() {
		wg.Wait()
		close(collect)
	}()

	var sitemapPages []page
	var blogPages []page
	for p := range collect {
		sitemapPages = append(sitemapPages, p)
		if p.File[0] == "blog" {
			blogPages = append(blogPages, p)
		}
	}
	sort.Slice(sitemapPages, func(i, j int) bool { return sitemapPages[i].CanonicalURL < sitemapPages[j].CanonicalURL })
	sort.Slice(blogPages, func(i, j int) bool { return blogPages[i].CanonicalURL < blogPages[j].CanonicalURL })

	// generate sitemap, blog index, atom feed

	return nil
}

type page struct {
	File         []string
	CanonicalURL string

	Title       []byte
	Description []byte
	Style       []byte
	Header      []byte
	Main        string
}

func (p page) String() string {
	return fmt.Sprintf("===== page: %s\nTitle: %s\nDescription: %s\nStyle: %s\nHeader: %s\nMain: %s\n\n", p.CanonicalURL, p.Title, p.Description, p.Style, p.Header, p.Main)
}

func (o *options) processPage(fn string) (page, error) {
	var p page
	b, err := ioutil.ReadFile(filepath.Join(o.src, fn))
	if err != nil {
		return p, fmt.Errorf("options.processPage: %s %w", fn, err)
	}
	bb := bytes.Split(b, []byte("---"))
	for _, b := range bb {
		if len(b) == 0 {
			continue
		}
		i := bytes.Index(b, []byte("\n"))
		switch string(bytes.TrimSpace(b[:i])) {
		case "title":
			p.Title = bytes.TrimSpace(b[i:])
		case "description":
			p.Description = bytes.TrimSpace(b[i:])
		case "style":
			p.Style = bytes.TrimSpace(b[i:])
		case "header":
			p.Header = bytes.TrimSpace(b[i:])
		case "main":

			p.Main = string(blackfriday.Run(b[i:], blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{HeadingLevelOffset: 3, Flags: blackfriday.CommonHTMLFlags}))))
			// p.Main = bytes.TrimSpace(b[i:])
		default:
			return p, fmt.Errorf("options.processPage: unknown section %s", bytes.TrimSpace(b[:i]))
		}
	}

	// TODO: write to files
	// f, err := openWrite("htmlfile")
	// if err != nil {
	//         return p, fmt.Errorf("options.processPage: %w", err)
	// }
	// defer f.Close()
	// o.templates.ExecuteTemplate(f, "htmlfile", p)
	// f, err = openWrite("ampfile")
	// if err != nil {
	//         return p, fmt.Errorf("options.processPage: %w", err)
	// }
	// defer f.Close()
	// o.templates.ExecuteTemplate(f, "ampfile", p)

	return p, nil
}

// func fileNames(fn, src, dst string) (dst, dstAMP, fullURL)

func openWrite(fn string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(fn), 0755)
	if err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}
	f, err := os.Create(fn)
	if err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}
	return f, nil
}
