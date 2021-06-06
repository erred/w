package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.seankhliao.com/w/v16/render"
)

func main() {
	matches, err := filepath.Glob("testdata/*.md")
	if err != nil {
		log.Fatalln("glob", err)
	}
	os.MkdirAll("out", 0o755)
	for _, m := range matches {
		func(m string) {
			fin, err := os.Open(m)
			if err != nil {
				log.Fatalln("open", m, err)
			}
			defer fin.Close()
			out := filepath.Join("out", strings.TrimSuffix(filepath.Base(m), ".md")+".html")
			fout, err := os.Create(out)
			if err != nil {
				log.Fatalln("create", out, err)
			}
			defer fout.Close()

			err = render.Render(&render.Options{Data: render.PageData{Compact: false}}, fout, fin)
			if err != nil {
				log.Fatalln("render", fin, err)
			}
		}(m)
	}
}
