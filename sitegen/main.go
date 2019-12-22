package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template"

	"github.com/russross/blackfriday/v2"
)

var (
	defaultSrc = "src"
	defaultDst = "dst"
	// defaultFonts     = "https://fonts.googleapis.com/css?family=Inconsolata:400,700|Lora:400,700&display=swap"
	// defaultAMPPrefux = "amp"
	defaultBaseURL = "https://seankhliao.com"
	defaultGAID    = "UA-114337586-1"
	defaultProject = "com-seankhliao"

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

	// err = o.deploy()
	// if err != nil {
	//                log.Println(err)
	//                os.Exit(1)
	//        }
}

type options struct {
	src string
	dst string
	// templateGlob string
	// fontURL      string
	// ampPrefix    string
	baseURL    string
	gaID       string
	gcpProject string

	templates *template.Template
}

func newOptions() *options {
	o := &options{}
	flag.StringVar(&o.src, "src", defaultSrc, "source directory")
	flag.StringVar(&o.dst, "dst", defaultDst, "output directory")
	// flag.StringVar(&o.templateGlob, defaultTemplates, "templates", "template glob")
	// flag.StringVar(&o.fontURL, "fonts", defaultFonts, "fonts to retreive")
	// flag.StringVar(&o.ampPrefix, "amp", defaultAMPPrefux, "amp path prefix")
	flag.StringVar(&o.baseURL, "host", defaultBaseURL, "url base")
	flag.StringVar(&o.gaID, "ga", defaultGAID, "google analytics ID")
	flag.StringVar(&o.gcpProject, "project", defaultProject, "GCP project (firebase)")

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

// func (o *options) getFonts() error {
// 	res, err := http.Get(o.fontURL)
// 	if err != nil {
// 		return fmt.Errorf("options.getFonts: %w", err)
// 	} else if res.StatusCode < 200 || res.StatusCode > 299 {
// 		return fmt.Errorf("options.getFonts: %d %s", res.StatusCode, res.Status)
// 	}
// 	defer res.Body.Close()
// 	buf := bytes.NewBufferString(`{{ define "fontcss" }}`)
// 	buf.ReadFrom(res.Body)
// 	buf.WriteString(`{{ end }}`)
//
// 	o.templates, err = o.templates.New("fontcss").Parse(buf.String())
// 	if err != nil {
// 		return fmt.Errorf("options.getFonts: %w", err)
// 	}
// 	return nil
// }

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
	ur, err := url.Parse(o.baseURL)
	if err != nil {
		return fmt.Errorf("options.parsePages: %w", err)
	}

	var wg sync.WaitGroup
	collect := make(chan []string)
	filepath.Walk(o.src, func(fp string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		} else if filepath.Ext(fp) == ".md" {
			wg.Add(1)
			go o.processPage(fp, collect, &wg)
		} else {
			wg.Add(1)
			go o.copyFile(fp, &wg)
		}
		return nil
	})
	go func() {
		wg.Wait()
		close(collect)
	}()

	var sitemapPages, blogPages []string
	for u := range collect {
		ur.Path = filepath.Join(u...)
		sitemapPages = append(sitemapPages, ur.String())
		if len(u) > 1 && u[0] == "blog" {
			blogPages = append(blogPages, u[len(u)-1])
		}
	}
	sort.Strings(sitemapPages)
	sort.Strings(blogPages)

	// generate sitemap, blog index, atom feed

	ioutil.WriteFile(filepath.Join(o.dst, "sitemap.txt"), []byte(strings.Join(sitemapPages, "\n")), 0644)
	return nil
}

type page struct {
	CanonicalURL string
	AMPURL       string

	Title       []byte
	Description []byte
	Style       []byte
	Header      []byte
	Main        string

	Date string // blogpost
}

func (p page) String() string {
	return fmt.Sprintf("===== page: %s\nTitle: %s\nDescription: %s\nStyle: %s\nHeader: %s\nMain: %s\n\n", p.CanonicalURL, p.Title, p.Description, p.Style, p.Header, p.Main)
}

// processPage takes a filepath.md (from current directory)
// and creates the cprresponding filepath.html and amp/filepath.html
// also sends the relative url path segments to collect
func (o *options) processPage(fp string, collect chan []string, done *sync.WaitGroup) {
	if done != nil {
		defer done.Done()
	}
	fps := strings.Split(fp, "/")
	var p page
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		// return fmt.Errorf("options.processPage: %s %w", fn, err)
		log.Printf("options.processPage: %s %v", fp, err)
		return
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
			p.Main = string(blackfriday.Run(b[i:], blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{HeadingLevelOffset: 2, Flags: blackfriday.CommonHTMLFlags}))))
		default:
			log.Printf("options.processPage: unknown section %s", bytes.TrimSpace(b[:i]))
			return
		}
	}

	// TODO: write to files
	fps[0], fps[len(fps)-1] = "amp", strings.TrimSuffix(fps[len(fps)-1], ".md")+".html"
	htmlpath := filepath.Join(fps[1:]...)
	f, err := openWrite(filepath.Join(o.dst, htmlpath))
	if err != nil {
		log.Printf("options.processPage: %v", err)
		return
	}
	defer f.Close()
	f.WriteString(p.Main)
	// o.templates.ExecuteTemplate(f, "htmlfile", p)
	collect <- fps[1:]

	amppath := filepath.Join(fps...)
	f, err = openWrite(filepath.Join(o.dst, amppath))
	if err != nil {
		log.Printf("options.processPage: %v", err)
		return
	}
	defer f.Close()
	f.WriteString(p.Main)
	// o.templates.ExecuteTemplate(f, "ampfile", p)
	collect <- fps
}

// func fileNames(fn, src, dst string) (dst, dstAMP, fullURL)

func (o *options) deploy() error {
	cmd := exec.Command("firebase", "-P", o.gcpProject, "deploy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("options.deploy: %w", err)
	}
	return nil
}

func (o *options) copyFile(fp string, done *sync.WaitGroup) {
	if done != nil {
		defer done.Done()
	}
	fps := strings.Split(fp, "/")
	f1, err := os.Open(fp)
	if err != nil {
		log.Println("options.parsePages: copy open f1", fp, err)
		return
	}
	defer f1.Close()

	fps[0] = o.dst
	f2, err := openWrite(filepath.Join(fps...))
	if err != nil {
		log.Println("options.parsePages: copy open f2", filepath.Join(fps...), err)
		return
	}
	defer f2.Close()
	io.Copy(f2, f1)
}

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
