package main

import (
	"flag"
	"os"

	"go.seankhliao.com/w/v16/process"
	"k8s.io/klog/v2/klogr"
)

func main() {
	var o process.Options
	var src, dst string
	flag.StringVar(&src, "src", "content", "source directory or file")
	flag.StringVar(&dst, "dst", "public", "destination directory or file")
	flag.StringVar(&o.Canonical, "url", "https://seankhliao.com", "base url for canonicalization")
	flag.StringVar(&o.GTMID, "gtm", "", "Google Tag Manager ID to enable analytics")
	flag.BoolVar(&o.Compact, "compact", false, "compact header")
	flag.BoolVar(&o.Raw, "raw", false, "skip markdown processing")
	flag.Parse()

	log := klogr.New()

	fi, err := os.Stat(src)
	if err != nil {
		log.Error(err, "stat", "src", fi.Name())
		os.Exit(1)
	}
	if fi.IsDir() {
		err = process.Dir(o, dst, src)
	} else {
		err = process.File(o, dst, src)
	}
	if err != nil {
		log.Error(err, "render")
		os.Exit(1)
	}
}
