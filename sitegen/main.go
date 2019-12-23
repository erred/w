package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/russross/blackfriday/v2"
)

func main() {
	o := newOptions()

	os.MkdirAll(o.dst, 0755)

	err := o.parseTemplates()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//
	// err = o.getFonts()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	//
	err = o.convertImgs()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	o.processPages()

	err = o.deploy()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type Page struct {
	URLCanonical string
	URLAMP       string
	AMP          bool
	GAID         string

	Title       string
	Description string
	Style       string
	Header      string
	Main        string

	Date  string     // blogpost
	Posts []BlogPost // blogindex
}

func (p *Page) setAMP() {
	p.AMP = true
}

type BlogPost struct {
	Date  string
	Title string
	URL   string
}

func (o options) processPages() error {
	var blogindex Page
	var sitemapPages []string
	var wg sync.WaitGroup

	sitemap, blog := make(chan string), make(chan BlogPost)
	filepath.Walk(o.src, func(fp string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		} else if filepath.Ext(fp) == ".md" {
			if filepath.Base(fp) == "blog.md" {
				_, blogindex, err = o.parsePage(fp)
				return nil
			}
			wg.Add(1)
			go o.processPage(fp, sitemap, blog, &wg)
		} else {
			wg.Add(1)
			go o.copyFile(fp, &wg)
		}
		return nil
	})
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(sitemap)
		close(blog)
		done <- struct{}{}
	}()

loop:
	for {
		select {
		case u := <-sitemap:
			sitemapPages = append(sitemapPages, u)
		case b := <-blog:
			blogindex.Posts = append(blogindex.Posts, b)
		case <-done:
			break loop
		}
	}
	sort.Strings(sitemapPages)
	sort.Slice(blogindex.Posts, func(i, j int) bool { return blogindex.Posts[i].URL > blogindex.Posts[j].URL })

	// generate sitemap, blog index, atom feed
	o.writeTemplate("/blog/index.html", "layout-blogindex", &blogindex)

	ioutil.WriteFile(filepath.Join(o.dst, "sitemap.txt"), []byte(strings.Join(sitemapPages, "\n")), 0644)

	return nil
}

// processPage takes a filepath from current directory
// and creates the cprresponding filepath.html and amp/filepath.html
// also sends the relative url path segments to collect
func (o options) processPage(fp string, sitemap chan string, blog chan BlogPost, done *sync.WaitGroup) {
	if done != nil {
		defer done.Done()
	}
	fps, p, err := o.parsePage(fp)
	if err != nil {
		log.Printf("options.processPage: %v", err)
		return
	}

	fps[len(fps)-1] = strings.TrimSuffix(fps[len(fps)-1], ".md") + ".html"
	htmlpath := filepath.Join(fps[1:]...)

	if fps[1] == "blog" {
		o.writeTemplate(htmlpath, "layout-blogpost", &p)
		blog <- BlogPost{
			Title: p.Title,
			Date:  p.Date,
			URL:   strings.TrimSuffix(fps[len(fps)-1], ".html"),
		}
	} else {
		o.writeTemplate(htmlpath, "layout-main", &p)
	}
	sitemap <- p.URLCanonical
	sitemap <- p.URLAMP
}

// parsePage takes a filepath from the current directory
// and returns a the path segments and a processed page
func (o options) parsePage(fp string) ([]string, Page, error) {
	fps, p := strings.Split(fp, "/"), Page{GAID: o.gaID}
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, p, fmt.Errorf("parsePage: %s %v", fp, err)
	}

	u, _ := url.Parse(o.baseURL)
	u.Path = filepath.Join(fps[1:]...)
	p.URLCanonical = normalizeURL(u.String())
	u.Path = filepath.Join(fps...)
	p.URLAMP = normalizeURL(u.String())

	bb := bytes.Split(b, []byte("---"))
	for _, b := range bb {
		if len(b) == 0 {
			continue
		}
		i := bytes.Index(b, []byte("\n"))
		switch string(bytes.TrimSpace(b[:i])) {
		case "title":
			p.Title = string(bytes.TrimSpace(b[i:]))
		case "description":
			p.Description = string(bytes.TrimSpace(b[i:]))
		case "style":
			p.Style = string(bytes.TrimSpace(b[i:]))
		case "header":
			p.Header = string(bytes.TrimSpace(b[i:]))
		case "main":
			p.Main = string(blackfriday.Run(
				b[i:],
				blackfriday.WithRenderer(
					blackfriday.NewHTMLRenderer(
						blackfriday.HTMLRendererParameters{
							HeadingLevelOffset: 2,
							Flags:              blackfriday.CommonHTMLFlags,
						},
					),
				),
			))
		default:
			return nil, p, fmt.Errorf("parsePage: unknown section %s", bytes.TrimSpace(b[:i]))
		}
	}

	if fps[1] == "blog" {
		p.Date = fps[len(fps)-1][:10]
	}
	return fps, p, nil
}
