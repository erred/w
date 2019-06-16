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
	"text/template"

	"github.com/russross/blackfriday/v2"
)

var (
	IgnoreExt = make(map[string]struct{})
	Src       = "src"
	Dst       = "dst"
	TmplExt   = ".gohtml"
	MdExt     = ".md"
	HtmlExt   = ".html"
	BaseURL   = "https://seankhliao.com"
)

func init() {
	var ignoreExt string
	flag.StringVar(&ignoreExt, "ignoreext", ".ico,.svg,.png,.jpg", "comma separated list of extensions to ignore")
	// name templates:
	// $basedir-post :ex /src/bog/template.gohtml => "blog-post"

	flag.Parse()
	for _, e := range strings.Split(ignoreExt, ",") {
		IgnoreExt[e] = struct{}{}
	}
}

func main() {
	p := NewProcessor()
	if err := filepath.Walk(Src, p.walker); err != nil {
		log.Fatal("main walker:", err)
	}
	p.Process()
}

type Processor struct {
	t *template.Template

	q       []*Page            // work queue
	indexes map[string][]*Page // generated index pages
	sitemap []string           // pages to include in sitemap
}

func NewProcessor() *Processor {
	return &Processor{
		t:       template.New(""),
		indexes: make(map[string][]*Page),
	}
}

// walker walks the src dir
// parses templates
// reads everything else into memory
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
		return nil
	}
	if ext == TmplExt {
		p.t, err = p.t.ParseFiles(fp)
		if err != nil {
			log.Printf("walker parse template %v: %v", fp, err)
		}
		return nil
	}

	page, err := NewPage(fp)
	if err != nil {
		log.Printf("walker newpage %v: %v", fp, err)
		return nil
	}
	if ext == MdExt {
		err = page.parseMD()
		if err != nil {
			log.Printf("walker parse markdown %v: %v", fp, err)
			return nil
		}
		p.indexes[filepath.Dir(fp)] = append(p.indexes[filepath.Dir(fp)], page)
	} else {
		p.q = append(p.q, page)
	}
	return nil
}

func (p *Processor) Process() {
	for _, page := range p.q {
		// fmt.Printf("processing %v %v\n", page.u, page.M)
		f, err := newFile(page.u.Dst())
		if err != nil {
			log.Printf("Process newFile %v: %v", page.u.Dst(), err)
			continue
		}
		if strings.HasSuffix(page.u.Dst(), HtmlExt) {
			p.sitemap = append(p.sitemap, page.u.Canonical())

			t, err := p.t.New("_page").Parse(string(page.b))
			if err != nil {
				log.Printf("Process parse %v as template: %v\n", page.u.Dst(), err)
				continue
			}
			// fmt.Println(t.DefinedTemplates())
			err = t.ExecuteTemplate(f, "_page", page.M)
			f.Close()
			if err != nil {
				log.Printf("Process execute %v: %v\n", page.u.Dst(), err)
			}
		} else {
			_, err = f.Write(page.b)
			f.Close()
			if err != nil {
				log.Printf("Process write %v: %v", page.u.Dst(), err)
			}
		}
	}
	type PL struct {
		Date  string
		Title string
		URL   string
	}
	for dir, pages := range p.indexes {
		idxpage := Page{
			u: NewURL(dir + "/index.html"),
			M: make(map[string]interface{}),
		}
		var posts []PL
		for _, page := range pages {
			// fmt.Printf("processing %v %v\n", page.u, page.M)
			posts = append(posts, PL{page.M["date"].(string), page.M["title"].(string), page.u.Relative()})
			f, err := newFile(page.u.Dst())
			if err != nil {
				log.Printf("Process newFile 2 %v: %v\n", page.u.Dst(), err)
				continue
			}
			err = p.t.ExecuteTemplate(f, filepath.Base(dir)+"-post", page.M)
			f.Close()
			if err != nil {
				log.Printf("Process execute 2 %v: %v\n", page.u.Dst(), err)
				continue
			}
		}
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Date > posts[j].Date
		})
		idxpage.M["posts"] = posts
		f, err := newFile(idxpage.u.Dst())
		if err != nil {
			log.Printf("Process newFile 3 %v: %v\n", idxpage.u.Dst(), err)
			continue
		}
		err = p.t.ExecuteTemplate(f, filepath.Base(dir)+"-index", idxpage.M)
		f.Close()
		if err != nil {
			log.Printf("Process execute 3 %v: %v\n", idxpage.u.Dst(), err)
			continue
		}

	}

	sort.Strings(p.sitemap)
	err := ioutil.WriteFile(filepath.Join(Dst, "sitemap.txt"), []byte(strings.Join(p.sitemap, "\n")), 0644)
	if err != nil {
		log.Printf("Process sitemap: %v\n", err)
	}

}

type URL struct {
	path    string
	dstPath string
}

// NewURL generates the url from a on disk file path
// /x/y.md              -> /x/y
// /index.html          -> /
// /x.html              -> /x
// /x/index.html        -> /x
// /x/y.html            -> /x/y
func NewURL(s string) URL {
	var u URL
	s = strings.TrimPrefix(s, Src)

	if strings.HasSuffix(s, MdExt) {
		s = strings.TrimSuffix(s, MdExt) + HtmlExt
	}
	u.dstPath = Dst + s

	s = strings.TrimSuffix(s, HtmlExt)
	s = strings.TrimSuffix(s, "/index")
	if s == "" {
		s = "/"
	}
	u.path = s

	return u
}

func (u URL) Dst() string {
	return u.dstPath
}

func (u URL) Canonical() string {
	return BaseURL + u.path
}

func (u URL) Relative() string {
	return u.path
}

type Page struct {
	u URL
	b []byte
	M map[string]interface{}
}

func NewPage(f string) (*Page, error) {
	u := NewURL(f)
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("NewPage readFile %v: %v", f, err)
	}
	m := map[string]interface{}{"url": u.Canonical()}
	return &Page{u, b, m}, nil
}

func (p *Page) parseMD() error {
	bb := bytes.SplitN(p.b, []byte("\n---\n"), 2)
	if len(bb) != 2 {
		return fmt.Errorf("parseMD split expected --- in %v\n", p.u.Dst())
	}

	for ln, line := range bytes.Split(bytes.TrimSpace(bb[0]), []byte("\n")) {
		l := bytes.SplitN(line, []byte("="), 2)
		if len(l) != 2 {
			return fmt.Errorf("parseMD metadata expected split by = in line %v in %v", ln, p.u.Dst())
		}
		p.M[string(bytes.TrimSpace(l[0]))] = string(bytes.TrimSpace(l[1]))
	}

	p.M["content"] = string(blackfriday.Run(bb[1], blackfriday.WithRenderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{HeadingLevelOffset: 3, Flags: blackfriday.CommonHTMLFlags}))))
	return nil
}

func newFile(fn string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(fn), 0755)
	if err != nil {
		return nil, fmt.Errorf("newFile mkdirall %v:%v", filepath.Dir(fn), err)
	}
	f, err := os.Create(fn)
	if err != nil {
		return nil, fmt.Errorf("newFile create %v: %v", fn, err)
	}
	return f, nil
}
