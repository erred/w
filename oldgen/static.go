package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type StaticData struct {
	Path string
	URL  string
	Desc string
}

// StaticOptions holds config needed for parsing static html pages
type StaticOptions struct {
	Src string
	Dst string
}

func NewStaticOptions(args []string) *StaticOptions {
	var o StaticOptions
	f := flag.NewFlagSet("static", flag.ExitOnError)
	f.StringVar(&o.Src, "src", "src", "source directory")
	f.StringVar(&o.Dst, "dst", "dst", "output directory")
	f.Parse(args)
	return &o
}

func (o *StaticOptions) Exec(opt *Options, pre, post *sync.WaitGroup) error {
	if pre != nil {
		pre.Wait()
	}
	if post != nil {
		defer post.Done()
	}
	var pages []StaticData
	filepath.Walk(o.Src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		subpath, _ := filepath.Rel(o.Src, path)

		b, err := ioutil.ReadFile(path)
		if err != nil {
			err = fmt.Errorf("StaticOptions.Exec walk read %q: %w", path, err)
			log.Println(err)
			return err
		}

		_, err = opt.T.New(subpath).Parse(string(b))
		if err != nil {
			err = fmt.Errorf("StaticOptions.Exec walk parse %q: %w", path, err)
			log.Println(err)
			return err
		}

		pages = append(pages, StaticData{
			subpath,
			filepath.Join(opt.host, canonicalURL(subpath)),
			"",
		})
		return nil
	})

	wg := &sync.WaitGroup{}
	for _, page := range pages {
		wg.Add(1)
		go func(page StaticData) {
			defer wg.Done()

			dfn := filepath.Join(o.Dst, strings.ReplaceAll(page.Path, ".gohtml", ".html"))
			err := writeTemplate(opt.T, page.Path, dfn, page)
			if err != nil {
				log.Printf("StaticOptions.Exec write %q: %v", dfn, err)
				return
			}
		}(page)
	}

	wg.Wait()
	return nil
}
