package main

import (
	"flag"
	"text/template"
)

var (
	defaultSrc     = "src"
	defaultDst     = "dst"
	defaultBaseURL = "https://seankhliao.com"
	defaultGAID    = "UA-114337586-1"
	defaultProject = "com-seankhliao"

	defaultImgArgs = [][]string{
		[]string{
			"-background", "none",
			"-density", "1200",
			"-resize", "1920x1080",
			"img/map.svg",
			"-write", "dst/map.png",
			"-write", "dst/map.webp",
			"dst/map.jpg",
		},
		[]string{
			"img/icon.tif",
			"-flatten",
			"(", "+clone", "-resize", "512x512", "-quality", "60", "-write", "dst/icon-512.png", "+delete", ")",
			"(", "+clone", "-resize", "192x192", "-quality", "60", "-write", "dst/icon-192.png", "+delete", ")",
			"(", "+clone", "-resize", "128x128", "-quality", "60", "-write", "dst/icon-128.png", "+delete", ")",
			"(", "+clone", "-resize", "64x64", "-quality", "60", "-write", "dst/icon-64.png", "+delete", ")",
			"(", "+clone", "-resize", "48x48", "-quality", "60", "-write", "dst/icon-48.png", "+delete", ")",
			"(", "+clone", "-resize", "32x32", "-quality", "60", "-write", "dst/icon-32.png", "+delete", ")",
			"(", "+clone", "-resize", "16x16", "-quality", "60", "-write", "dst/icon-16.png", "+delete", ")",
			"-resize", "32x32", "dst/favicon.ico",
		},
	}
)

type options struct {
	src        string
	dst        string
	baseURL    string
	gaID       string
	gcpProject string

	templates *template.Template
}

func newOptions() *options {
	o := &options{}
	flag.StringVar(&o.src, "src", defaultSrc, "source directory")
	flag.StringVar(&o.dst, "dst", defaultDst, "output directory")
	flag.StringVar(&o.baseURL, "host", defaultBaseURL, "url base")
	flag.StringVar(&o.gaID, "ga", defaultGAID, "google analytics ID")
	flag.StringVar(&o.gcpProject, "project", defaultProject, "GCP project (firebase)")

	flag.Parse()
	return o
}
