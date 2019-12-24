package main

import (
	"flag"
	"fmt"
	"os"
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

	// SXG
	SXG         bool
	certURL     string
	validityURL string
	certPath    string
	privPath    string

	templates *template.Template
}

func newOptions() (*options, error) {
	o := &options{}
	flag.StringVar(&o.src, "src", defaultSrc, "source directory")
	flag.StringVar(&o.dst, "dst", defaultDst, "output directory")
	flag.StringVar(&o.baseURL, "host", defaultBaseURL, "url base")
	flag.StringVar(&o.gaID, "ga", defaultGAID, "google analytics ID")
	flag.StringVar(&o.gcpProject, "project", defaultProject, "GCP project (firebase)")

	flag.BoolVar(&o.SXG, "sxg", false, "enable HTTP SXG")
	flag.StringVar(&o.certURL, "certURL", defaultBaseURL+"/cert.cbor", "url to find the SXG signing certificate (CBOR)")
	flag.StringVar(&o.certPath, "certPath", "/var/secrets/cbor.pem", "path to find the SXG signing certificate (CBOR)")
	flag.StringVar(&o.validityURL, "validityURL", defaultBaseURL+"/resource.validity.msg", "TODO: find out what this is")
	flag.StringVar(&o.privPath, "privPath", "/var/secrets/cbor.key", "path to dind the SXG signing key")

	flag.Parse()

	o.templates = template.New("")
	for name, tmpl := range rawTemplates {
		o.templates = template.Must(o.templates.New(name).Parse(tmpl))
	}

	err := os.MkdirAll(o.dst, 0755)
	if err != nil {
		return nil, fmt.Errorf("newOptions: create dst dir %s: %w", o.dst, err)
	}

	return o, nil
}
