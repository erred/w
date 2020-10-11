package main

import (
	"flag"
	"log"
	"os"

	"go.seankhliao.com/com-seankhliao/v13/internal/render"
)

func main() {
	var src, dst, baseURL string
	var embedStyle, disableAanalytics bool
	flag.StringVar(&src, "src", "site", "source directory or file")
	flag.StringVar(&dst, "dst", "public", "destination directory or file")
	flag.StringVar(&baseURL, "url", "https://seankhliao.com", "base url for canonicalization")
	flag.BoolVar(&embedStyle, "embedStyle", false, "embed stylesheets")
	flag.BoolVar(&disableAanalytics, "disableAanalytics", false, "disable google analytics")
	flag.Parse()

	fi, err := os.Stat(src)
	if err != nil {
		log.Fatalf("stat src=%s err: %v", fi, err)
	}
	if fi.IsDir() {
		err = render.ProcessDir(src, dst, baseURL, disableAanalytics, embedStyle)
	} else {
		_, err = render.ProcessFile(src, dst, baseURL, disableAanalytics, embedStyle)
	}
	if err != nil {
		log.Fatalf("render: %v", err)
	}
}
