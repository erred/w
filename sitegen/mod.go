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

// ModData holds old names (github user/repo)
// and new module names
type ModData struct {
	Old string
	New string
}

// ModOptions holds config needed to parse the gomod.txt file
// where each line is "user/repo module(without hostname)"
// and create the pages for redirects
// external dependencies:
//      template called "gomod"
//
type ModOptions struct {
	Src string
	Dst string
}

func NewModOptions(args []string) *ModOptions {
	var o ModOptions
	f := flag.NewFlagSet("mod", flag.ExitOnError)
	f.StringVar(&o.Src, "src", "gomod.txt", "source file")
	f.StringVar(&o.Dst, "dst", "dst", "output directory")
	f.Parse(args)
	return &o
}
func (o *ModOptions) Exec(opt *Options) error {
	b, err := ioutil.ReadFile(o.Src)
	if err != nil {
		return fmt.Errorf("ModOptions.Exec read file %q: %w", o.Src, err)
	}

	var mds []ModData
	for i, line := range bytes.Split(b, []byte("\n")) {
		fields := bytes.Fields(line)
		if len(fields) == 0 {
			continue
		} else if len(fields) != 2 {
			log.Printf("parsing %q: line %d expected 2 fields, got %d\n", o.Src, i, len(fields))
			continue
		}
		mds = append(mds, ModData{string(fields[0]), filepath.Join(opt.host, string(fields[1]))})
	}

	wg := &sync.WaitGroup{}
	for _, md := range mds {
		wg.Add(1)
		go func(md ModData) {
			defer wg.Done()
			var b bytes.Buffer
			err := opt.T.ExecuteTemplate(&b, "gomod", md)
			if err != nil {
				log.Printf("ModOptions.Exec template for %q: %w", md, err)
				return
			}
			dst := filepath.Join(o.Dst, strings.TrimPrefix(md.New, opt.host)+".html")
			os.MkdirAll(filepath.Dir(dst), 0755)
			err = ioutil.WriteFile(dst, b.Bytes(), 0644)
			if err != nil {
				log.Printf("ModOptions.Exec write %q: %w", dst, err)
			}
		}(md)
	}
	wg.Wait()
	return nil
}
