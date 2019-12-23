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
	// err = o.convertImgs()
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }

	err = o.processPages()
	if err != nil {
		log.Println(err)
	}

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
	var blogPosts []BlogPost

	err := filepath.Walk(o.src, func(fp string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		} else if filepath.Ext(fp) == ".md" {
			if strings.HasSuffix(fp, "blog/index.md") {
				_, blogindex, err = o.parsePage(fp)
				return nil
			}
			urls, bp, err := o.processPage(fp)
			if err != nil {
				return fmt.Errorf("options.processPages: %w", err)
			}
			sitemapPages = append(sitemapPages, urls...)
			if bp != nil {
				blogPosts = append(blogPosts, *bp)
			}
		} else {
			o.copyFile(fp, nil)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("options.processPages: %w", err)
	}

	sort.Strings(sitemapPages)
	sort.Slice(blogPosts, func(i, j int) bool { return blogPosts[i].URL > blogPosts[j].URL })
	blogindex.Posts = blogPosts

	// generate sitemap, blog index, atom feed
	err = o.writeTemplate("/blog/index.html", "layout-blogindex", &blogindex)
	if err != nil {
		return fmt.Errorf("options.processPages: %w", err)
	}
	err = ioutil.WriteFile(filepath.Join(o.dst, "sitemap.txt"), []byte(strings.Join(sitemapPages, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("options.processPages: %w", err)
	}

	return nil
}

// processPage takes a filepath from current directory
// and creates the cprresponding filepath.html and amp/filepath.html
// also sends the relative url path segments to collect
func (o options) processPage(fp string) (urls []string, bp *BlogPost, err error) {
	fps, p, err := o.parsePage(fp)
	if err != nil {
		return nil, nil, fmt.Errorf("options.processPage: %w", err)
	}

	htmlpath := filepath.Join(fps[1:]...) + ".html"
	if fps[1] == "blog" {
		err = o.writeTemplate(htmlpath, "layout-blogpost", &p)
		if err != nil {
			return nil, nil, fmt.Errorf("options.processPage: %w", err)
		}
		if p.Title == "" {
			fmt.Println(fp)
		}
		bp = &BlogPost{
			Title: p.Title,
			Date:  p.Date,
			URL:   strings.TrimSuffix(fps[len(fps)-1], ".html"),
		}
	} else {
		err = o.writeTemplate(htmlpath, "layout-main", &p)
		if err != nil {
			return nil, nil, fmt.Errorf("options.processPage: %w", err)
		}
	}
	return []string{p.URLCanonical, p.URLAMP}, bp, nil
}

// parsePage takes a filepath from the current directory
// and returns a the path segments and a processed page
func (o options) parsePage(fp string) ([]string, Page, error) {
	fps, p := strings.Split(fp, "/"), Page{GAID: o.gaID}
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, p, fmt.Errorf("parsePage: %s %v", fp, err)
	}
	fps[0], fps[len(fps)-1] = "amp", strings.TrimSuffix(fps[len(fps)-1], ".md")
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

	if fps[1] == "blog" && fps[2] != "index" {
		p.Date = fps[len(fps)-1][:10]
	}
	return fps, p, nil
}
