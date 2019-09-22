package main

import (
	"bytes"
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

func (o *StaticOptions) Exec(opt *Options) error {
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

		t := opt.T.New(subpath)
		_, err = t.Parse(string(b))
		if err != nil {
			err = fmt.Errorf("StaticOptions.Exec walk parse %q: %w", path, err)
			log.Println(err)
			return err
		}

		pages = append(pages, StaticData{subpath})
		return nil
	})

	wg := &sync.WaitGroup{}
	for _, page := range pages {
		wg.Add(1)
		go func(page StaticData) {
			defer wg.Done()

			var b bytes.Buffer
			err := opt.T.ExecuteTemplate(&b, page.Path, page)
			if err != nil {
				log.Printf("StaticOptions.Exec exec template %q: %w", page.Path, err)
				return
			}

			dst := filepath.Join(o.Dst, strings.ReplaceAll(page.Path, ".gohtml", ".html"))
			err = ioutil.WriteFile(dst, b.Bytes(), 0644)
			if err != nil {
				log.Printf("StaticOptions.Exec write %q: %w", dst, err)
				return
			}
		}(page)
	}

	wg.Wait()
	return nil
}
