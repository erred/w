package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

func main() {
	o, err := NewOptions("seankhliao.com", "templates")
	if err != nil {
		log.Fatal("initialize options: ", err)
	}

	var wg, remote, sitemap sync.WaitGroup

	wg.Add(1)
	go NewImgOptions(nil).Exec(o, nil, &wg)

	remote.Add(1)
	go NewRemoteOptions(nil).Exec(o, nil, &remote)

	sitemap.Add(1)
	go NewModOptions(nil).Exec(o, nil, &sitemap)

	sitemap.Add(1)
	go NewBlogOptions(nil).Exec(o, &remote, &sitemap)

	sitemap.Add(1)
	go NewStaticOptions(nil).Exec(o, &remote, &sitemap)

	wg.Add(1)
	go NewSitemapOptions(nil).Exec(o, &sitemap, &wg)

	wg.Wait()
}

// Options contains global options
// such as the template dir
// and host
type Options struct {
	host string             // hostname
	t    string             // template file or dir (flat)
	T    *template.Template // parsed templates
}

// NewOptions parses global options and decides on subcommand
// remember to register here
func NewOptions(host, t string) (*Options, error) {
	o := &Options{
		host: host,
		t:    filepath.Join(t, "*.gohtml"),
	}

	var err error
	o.T, err = template.ParseGlob(o.t)
	if err != nil {
		return nil, fmt.Errorf("NewOptions parse file %q: %w", o.t, err)
	}

	return o, nil
}

func canonicalURL(subpath string) string {
	subpath = strings.TrimSuffix(subpath, ".html")
	subpath = strings.TrimSuffix(subpath, "index")
	subpath = strings.TrimSuffix(subpath, "/")
	if subpath == "" {
		subpath = "/"
	}
	return subpath
}

func writeTemplate(t *template.Template, tname, fname string, data interface{}) error {
	os.MkdirAll(filepath.Dir(fname), 0755)
	f, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("writeTemplate open %v: %w", fname, err)
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, tname, data)
	if err != nil {
		return fmt.Errorf("writeTemplate execute %v: %w", tname, err)
	}
	return nil
}
