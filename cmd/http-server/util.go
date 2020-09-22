package main

import (
	"net/http"
	"os"
	"strings"
)

func exists(p string) bool {
	fi, err := os.Stat(p)
	if err != nil || fi.IsDir() {
		return false
	}
	return true
}

func canonical(p string) string {
	p = strings.TrimSuffix(strings.TrimSuffix(p, ".html"), "index")
	if p[len(p)-1] != '/' {
		p = p + "/"
	}
	return p
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("strict-transport-security", `max-age=63072000; preload`)
	w.Header().Set("referrer-policy", "strict-origin-when-cross-origin")
	w.Header().Set("report-to", `{"group": "csp-endpoint", "max_age": 10886400, "endpoints": [{"url":"https://statslogger.seankhliao.com/json"}]}`)
	w.Header().Set("content-security-policy", `default-src 'self'; upgrade-insecure-requests; connect-src https://statslogger.seankhliao.com https://www.google-analytics.com; font-src https://seankhliao.com; img-src *; object-src 'none'; script-src-elem 'nonce-deadbeef2' 'nonce-deadbeef3' 'nonce-deadbeef4' https://unpkg.com https://www.google-analytics.com https://ssl.google-analytics.com https://www.googletagmanager.com; sandbox allow-scripts; style-src-elem 'nonce-deadbeef1' https://seankhliao.com; report-to csp-endpoint; report-uri https://statslogger.seankhliao.com/json`)
	w.Header().Set("feature-policy", `accelerometer 'none'; autoplay 'none'; camera 'none'; document-domain 'none'; encrypted-media 'none'; fullscreen 'none'; geolocation 'none'; gyroscope 'none'; magnetometer 'none'; microphone 'none'; midi 'none'; payment 'none'; picture-in-picture 'none'; sync-xhr 'none'; usb 'none'; xr-spatial-tracking 'none'`)
}
