package main

import (
	"bytes"
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
	src = "src"
	dst = "dst"
	ext = map[string]string{
		".md": ".html",
	}
)

func main() {
	w := NewW()
	w.Start()
}

type W struct {
	wg    sync.WaitGroup
	t     *template.Template
	posts chan Post
}

func NewW() *W {
	t := template.Must(template.New("head").Parse(HeadTemplate))
	t = template.Must(t.New("foot").Parse(FootTemplate))
	t = template.Must(t.New("post").Parse(PostTemplate))
	t = template.Must(t.New("index").Parse(IndexTemplate))

	return &W{
		posts: make(chan Post, 32),
		t:     t,
	}
}

func (w *W) Start() {
	go func() {
		filepath.Walk(src, w.walker)
		w.wg.Wait()
		close(w.posts)
	}()
	w.writeIndex()
}

func (w *W) walker(path string, info os.FileInfo, err error) error {
	if err != nil || info.IsDir() {
		return nil
	}
	switch filepath.Ext(path) {
	case ".md":
		w.Convert(path)
	default:
		w.Copy(path)
	}
	return nil
}

func (w *W) writeIndex() {
	var ps []Post
	for p := range w.posts {
		ps = append(ps, p)
	}
	sort.Sort(Posts(ps))

	nfn := filepath.Join(dst, "blog/index.html")
	d, err := create(nfn)
	if err != nil {
		log.Printf("Copy error opening %v: %v\n", nfn, err)
		return
	}
	defer d.Close()

	idx := Index{Posts: ps, Description: "blog of seankhliao"}

	err = w.t.ExecuteTemplate(d, "index", idx)
	if err != nil {
		log.Printf("error executing template %v for %v: %v\n", "index", nfn, err)
		return
	}
}

func (w *W) Copy(fn string) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		s, err := os.Open(fn)
		if err != nil {
			log.Printf("Copy error opening %v: %v\n", fn, err)
			return
		}
		defer s.Close()

		d, err := create(rename(fn))
		if err != nil {
			log.Printf("Copy error opening %v: %v\n", rename(fn), err)
			return
		}
		defer d.Close()

		_, err = io.Copy(d, s)
		if err != nil {
			log.Printf("Copy error copy %v: %v", fn, err)
		}
	}()
}

func (w *W) Convert(fn string) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		b, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Printf("Convert error reading file %v: %v\n", fn, err)
			return
		}

		nfn := rename(fn)
		p := convert(nfn, b)

		d, err := create(nfn)
		if err != nil {
			log.Printf("Copy error opening %v: %v\n", nfn, err)
			return
		}
		defer d.Close()

		err = w.t.ExecuteTemplate(d, "post", p)
		if err != nil {
			log.Printf("error executing template %v for %v: %v\n", "post", nfn, err)
			return
		}

		w.posts <- p
	}()
}

func create(fn string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(fn), 0744)
	if err != nil {
		log.Printf("create error mkdirall %v: %v\n", filepath.Dir(fn), err)
		return nil, err
	}

	f, err := os.Create(fn)
	if err != nil {
		log.Printf("create error opening %v: %v\n", fn, err)
	}
	return f, err
}

func urlname(p string) string {
	return filepath.Join(filepath.SplitList(strings.TrimSuffix(p, filepath.Ext(p)))[1:]...)
}

func rename(p string) string {
	e := filepath.Ext(p)
	p = strings.TrimSuffix(p, e) + ext[e]
	p = dst + strings.TrimPrefix(p, src)
	return p
}

func convert(fn string, b []byte) Post {
	bb := bytes.SplitN(b, []byte("\n---\n"), 2)
	if len(bb) < 2 {
		log.Printf("convert --- not found in %v\n", fn)
	}

	title, date, desc := parseMeta(bb[0])

	p := Post{
		URL:         urlname(fn),
		Title:       title,
		Date:        date,
		Description: desc,
		Content:     string(blackfriday.Run(bb[1])),
	}
	return p
}

func parseMeta(b []byte) (title, date, desc string) {
	for _, b := range bytes.Split(bytes.TrimSpace(b), []byte("\n")) {
		bb := bytes.SplitN(b, []byte("="), 2)
		if len(bb) < 2 {
			log.Printf("parseMeta expected split by =\n")
			continue
		}
		v := string(bytes.TrimSpace(bb[1]))
		switch string(bytes.TrimSpace(bb[0])) {
		case "title":
			title = v
		case "date":
			date = v
		case "description", "desc":
			desc = v
		}
	}
	return title, date, desc
}

type Post struct {
	Title       string
	URL         string
	Description string
	Date        string
	Content     string
}

type Index struct {
	Posts []Post

	// dummy for template
	Title       string
	Description string
	URL         string
}

// newer first
type Posts []Post

func (p Posts) Less(i, j int) bool { return p[i].Date > p[j].Date }
func (p Posts) Len() int           { return len(p) }
func (p Posts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
