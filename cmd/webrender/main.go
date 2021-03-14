package main

import (
	"flag"
	"os"

	"go.seankhliao.com/w/v15/internal/render"
	"k8s.io/klog/v2/klogr"
)

func main() {
	var src, dst, baseURL string
	var embedStyle, disableAanalytics bool
	flag.StringVar(&src, "src", "content", "source directory or file")
	flag.StringVar(&dst, "dst", "public", "destination directory or file")
	flag.StringVar(&baseURL, "url", "https://seankhliao.com", "base url for canonicalization")
	flag.BoolVar(&embedStyle, "embedstyle", false, "embed stylesheets")
	flag.BoolVar(&disableAanalytics, "disableanalytics", false, "disable google analytics")
	flag.Parse()

	log := klogr.New()

	fi, err := os.Stat(src)
	if err != nil {
		log.Error(err, "stat", "src", fi.Name())
		os.Exit(1)
	}
	if fi.IsDir() {
		err = render.ProcessDir(src, dst, baseURL, disableAanalytics, embedStyle)
	} else {
		_, err = render.ProcessFile(src, dst, baseURL, disableAanalytics, embedStyle)
	}
	if err != nil {
		log.Error(err, "render")
		os.Exit(1)
	}
}
