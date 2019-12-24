package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WICG/webpackage/go/signedexchange"
	"github.com/WICG/webpackage/go/signedexchange/version"
)

// signExchanges signs all resources as SXG
// default settings:
//      method = GET
//      MIRecordSize = 4096
// common for all requests:
//      certificate
//      privatekey
//      certurl
//      validtyurl
//      expire
// per request
//      payload
//      uri
//      content-type (header)
func (o *options) signExchanges() error {
	if !o.SXG {
		return nil
	}
	signer, err := o.createSigner()
	if err != nil {
		return fmt.Errorf("options.signExchanges: %w", err)
	}
	err = filepath.Walk(o.dst, func(fp string, fi os.FileInfo, err error) error {
		if fi.IsDir() || filepath.Ext(fp) == ".sxg" {
			return nil
		}
		uri, resHeader, err := signGetURICT(o.baseURL, fp, o.dst)
		if err != nil {
			return fmt.Errorf("walker: %w", err)
		}

		payload, err := ioutil.ReadFile(fp)
		if err != nil {
			return fmt.Errorf("walker: %s %w", fp, err)
		}
		e := signedexchange.NewExchange(version.Version1b3, uri, http.MethodGet, http.Header{}, http.StatusOK, resHeader, payload)
		err = e.MiEncodePayload(4096)
		if err != nil {
			return fmt.Errorf("walker: miencode %w", err)
		}
		err = e.AddSignatureHeader(signer)
		if err != nil {
			return fmt.Errorf("walker: addsig %w", err)
		}
		f, err := os.Create(fp + ".sxg")
		if err != nil {
			return fmt.Errorf("walker: open outfile %w", err)
		}
		defer f.Close()
		err = e.Write(f)
		if err != nil {
			return fmt.Errorf("walker: write %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("options.signExchanges: walk: %w", err)
	}
	return nil
}

func (o *options) createSigner() (*signedexchange.Signer, error) {
	certURL, err := url.Parse(o.certURL)
	if err != nil {
		return nil, fmt.Errorf("options.createSigner: parse certURL: %w", err)
	}
	validityURL, err := url.Parse(o.validityURL)
	if err != nil {
		return nil, fmt.Errorf("options.createSigner: parse validityURL: %w", err)
	}

	certsb, err := ioutil.ReadFile(o.certPath)
	if err != nil {
		return nil, fmt.Errorf("options.createSigner: read certPath: %w", err)
	}
	certs, err := signedexchange.ParseCertificates(certsb)
	if err != nil {
		return nil, fmt.Errorf("options.createSigner: parse certs: %w", err)
	}

	pkeyb, err := ioutil.ReadFile(o.privPath)
	if err != nil {
		return nil, fmt.Errorf("options.createSigner: read privPath: %w", err)
	}
	privKey, err := signedexchange.ParsePrivateKey(pkeyb)
	if err != nil {
		return nil, fmt.Errorf("options.createSigner: parse privKey: %w", err)
	}

	return &signedexchange.Signer{
		Date: time.Now(),
		// Max, see section 3.5
		// https://wicg.github.io/webpackage/draft-yasskin-http-origin-signed-responses.html#name-signature-validity
		Expires:     time.Now().Add(7 * 24 * time.Hour),
		Certs:       certs,
		CertUrl:     certURL,
		ValidityUrl: validityURL,
		PrivKey:     privKey,
	}, nil
}

func signGetURICT(baseURL, fp, dst string) (string, http.Header, error) {
	ct := "text/html; charset=utf-8"
	resHeader := http.Header{}
	switch filepath.Ext(fp) {
	case ".html":
		ct = ct
	case ".png":
		ct = "image/png"
	case ".webp":
		ct = "image/webp"
	case ".jpg":
		ct = "image/jpeg"
	case ".json":
		ct = "application/json"
	case ".cbor":
		ct = "application/cert-chain+cbor"
	case ".ico":
		ct = "image/x-icon"
	case ".atom":
		ct = "application/rss+xml"
	case ".txt":
		ct = "text/plain"
	default:
		return "", nil, fmt.Errorf("options.signExchanges: unknown content type for %s", fp)

	}
	resHeader.Add("content-type", ct)
	u, _ := url.Parse(baseURL)
	u.Path = strings.TrimPrefix(fp, dst)
	uri := normalizeURL(u.String())
	return uri, resHeader, nil
}
