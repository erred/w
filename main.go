package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template"

	"github.com/russross/blackfriday/v2"
)

var (
	IgnoreExt = make(map[string]struct{})
	Src, Dst  string
	TmplExt   string
	MdExt     = ".md"
	HtmlExt   = ".html"
	BaseURL   = "https://seankhliao.com"
)

func init() {
	var ignoreExt string
	flag.StringVar(&ignoreExt, "ignoreext", ".ico,.svg,.png,.jpg", "comma separated list of extensions to ignore")
	flag.StringVar(&Src, "src", "src", "source directory")
	flag.StringVar(&Dst, "dst", "dst", "output directory")
	// name templates:
	// $basedir-post :ex /src/bog/template.gohtml => "blog-post"
	flag.StringVar(&TmplExt, "tmplext", ".gohtml", "extension for html templates")

	flag.Parse()
	for _, e := range strings.Split(ignoreExt, ",") {
		IgnoreExt[e] = struct{}{}
	}
}

func main() {
	p := NewProcessor()
	p.Walk()
	p.Process()
	p.GenSitemap()
}

type Processor struct {
	t             *template.Template
	paths         []string
	templatePaths []string
	renderPaths   []string
	copyPaths     []string
	mdPaths       map[string][]string
}

func NewProcessor() *Processor {
	return &Processor{
		mdPaths: make(map[string][]string),
	}
}

func (p *Processor) Walk() {
	filepath.Walk(Src, p.walker)
}
func (p *Processor) Process() error {
	var wg sync.WaitGroup
	var err error

	p.t, err = template.ParseFiles(p.templatePaths...)
	if err != nil {
		log.Printf("Process parseFiles %v\n", err)
		return err
	}

	defer wg.Wait()

	for _, fp := range p.renderPaths {
		wg.Add(1)
		go func(fp string) {
			defer wg.Done()
			p.renderFile(fp)
		}(fp)
	}

	for _, fp := range p.copyPaths {
		wg.Add(1)
		go func(fp string) {
			defer wg.Done()
			copyFile(fp)
		}(fp)
	}

	for dir, fps := range p.mdPaths {
		wg.Add(1)
		go func(dir string, fps []string) {
			defer wg.Done()
			p.md2html(dir, fps)
		}(dir, fps)
	}

	return nil
}

func (p *Processor) walker(fp string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("walker called with %v\n", err)
		return nil
	}
	if info.IsDir() {
		if strings.HasSuffix(fp, "-src") {
			return filepath.SkipDir
		}
		return nil
	}

	ext := filepath.Ext(fp)
	if _, ok := IgnoreExt[ext]; ok {
		// noop
	} else if ext == TmplExt {
		p.templatePaths = append(p.templatePaths, fp)
	} else if ext == MdExt {
		dir := filepath.Dir(fp)
		p.mdPaths[dir] = append(p.mdPaths[dir], fp)
	} else if ext == HtmlExt {
		p.renderPaths = append(p.renderPaths, fp)
	} else {
		p.copyPaths = append(p.copyPaths, fp)
	}

	return nil
}

func (p *Processor) GenSitemap() {
	filepath.Walk(Dst, func(fp string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("GenSitemap called with %v\n", err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(fp) != HtmlExt {
			return nil
		}
		if relativeURL(fp) == "/404" {
			return nil
		}
		p.paths = append(p.paths, canonicalURL(fp))
		return nil
	})

	sort.Strings(p.paths)
	err := ioutil.WriteFile(filepath.Join(Dst, "sitemap.txt"), []byte(strings.Join(p.paths, "\n")), 0644)
	if err != nil {
		log.Printf("GenSitemap write sitemap: %v\n", err)
	}
}

func (p *Processor) md2html(dir string, fps []string) {
	var posts []Post
	bd := filepath.Base(dir)
	posttmpl := bd + "-post"
	idxtmpl := bd + "-index"

	for _, fp := range fps {
		posts = append(posts, p.parsePost(posttmpl, fp))
	}

	sort.Sort(Posts(posts))

	idx := map[string]interface{}{
		"Posts": posts,
		"URL":   relativeURL(dir),
	}

	nfn := filepath.Join(src2dst(dir), "index.html")
	f, err := newFile(nfn)
	if err != nil {
		log.Printf("md2html newFile %v: %v\n", nfn, err)
		return
	}
	defer f.Close()

	err = p.t.ExecuteTemplate(f, idxtmpl, idx)
	if err != nil {
		log.Printf("md2html exec template %v for %v: %v", idxtmpl, nfn, err)
		return
	}

}

func (p *Processor) parsePost(tmpl, fp string) Post {
	var pt Post

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Printf("parsePost readfile %v: %v\n", fp, err)
		return pt
	}

	nfn := strings.TrimSuffix(src2dst(fp), MdExt) + ".html"
	pt.URL = relativeURL(nfn)

	bb := bytes.SplitN(b, []byte("\n---\n"), 2)
	if len(bb) != 2 {
		log.Printf("parsePost split expected --- in %v\n", fp)
		return pt
	}
	for ln, line := range bytes.Split(bytes.TrimSpace(bb[0]), []byte("\n")) {
		l := bytes.SplitN(line, []byte("="), 2)
		if len(l) != 2 {
			log.Printf("parsePost parse header %v line %v expected split by = \n", fp, ln)
			return pt
		}
		v := string(bytes.TrimSpace(l[1]))
		switch string(bytes.TrimSpace(l[0])) {
		case "title":
			pt.Title = v
		case "date":
			pt.Date = v
		case "desc", "description":
			pt.Desc = v
		default:
			log.Printf("parsePost parse header %v line %v unkown kv: %v\n", fp, ln, l)
		}
	}

	pt.Content = string(blackfriday.Run(bytes.TrimSpace(bb[1])))

	f, err := newFile(nfn)
	if err != nil {
		log.Printf("parsePost newFile %v: %v\n", nfn, err)
		return pt
	}
	defer f.Close()

	err = p.t.ExecuteTemplate(f, tmpl, pt)
	if err != nil {
		log.Printf("parsePost exec template %v for %v: %v", tmpl, fp, err)
		return pt
	}
	return pt
}

func (p *Processor) renderFile(fn string) {
	nfn := src2dst(fn)
	f, err := newFile(nfn)
	if err != nil {
		log.Printf("renderFile newFile %v: %v\n", nfn, err)
		return
	}
	t, err := p.t.ParseFiles(fn)
	if err != nil {
		log.Printf("renderFile parse %v as template: %v\n", fn, err)
		return
	}
	t.ExecuteTemplate(f, filepath.Base(fn), nil)
	if err != nil {
		log.Printf("renderFile execute %v: %v\n", fn, err)
	}
}

func copyFile(fn string) {
	nfn := src2dst(fn)
	f, err := newFile(nfn)
	if err != nil {
		log.Printf("copyFile newFile %v: %v\n", nfn, err)
		return
	}
	fo, err := os.Open(fn)
	if err != nil {
		log.Printf("copyFile Open %v: %v\n", fn, err)
		return
	}
	defer fo.Close()

	_, err = io.Copy(f, fo)
	if err != nil {
		log.Printf("copyFile copy %v tp %v: %v", fn, nfn, err)
	}
}

func relativeURL(fp string) string {
	fp = strings.TrimPrefix(fp, Src)
	fp = strings.TrimPrefix(fp, Dst)
	fp = strings.TrimSuffix(fp, HtmlExt)
	fp = strings.TrimSuffix(fp, "index")
	if len(fp) != 1 {
		fp = strings.TrimSuffix(fp, "/")
	}
	return fp
}
func canonicalURL(fp string) string {
	fp = BaseURL + relativeURL(fp)
	return strings.TrimSuffix(fp, "/")
}

func newFile(fn string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(fn), 0755)
	if err != nil {
		log.Printf("newFile mkdirall %v:%v\n", filepath.Dir(fn), err)
	}
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("newFile create %v: %v\n", fn, err)
		return nil, err
	}
	return f, nil
}
func src2dst(f string) string {
	return filepath.Join(Dst, strings.TrimPrefix(f, Src))
}

type Post struct {
	Title   string
	URL     string
	Desc    string
	Date    string
	Content string
}

// newer first
type Posts []Post

func (p Posts) Less(i, j int) bool { return p[i].Date > p[j].Date }
func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
