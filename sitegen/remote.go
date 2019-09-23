package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// RemoteOptions holds config for fetching remote content
// we want to include statically
// such as Google Font's css file,
// Creates:
//      template named "fontcss"
type RemoteOptions struct {
	FontURL string
	Dst     string
}

func NewRemoteOptions(args []string) *RemoteOptions {
	u := "https://fonts.googleapis.com/css?family=Inconsolata:400,700|Lora:400,700&display=swap"
	var o RemoteOptions
	f := flag.NewFlagSet("remote", flag.ExitOnError)
	f.StringVar(&o.FontURL, "fonturl", u, "google fonts url")
	f.StringVar(&o.Dst, "dst", "templates/fontcss.gohtml", "output file")
	f.Parse(args)
	return &o
}
func (o *RemoteOptions) Exec(opt *Options) error {
	r, err := http.Get(o.FontURL)
	if err != nil {
		return fmt.Errorf("RemoteOptions.Exec get %q: %w", o.FontURL, err)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("RemoteOptions.Exec read resp: %w", err)
	}

	fi, err := os.Stat(o.Dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("RemoteOptions.Exec stat %q: %w", o.Dst, err)
		}
	} else {
		// exists, remove
		if fi.IsDir() {
			return fmt.Errorf("RemoteOptions.Exec %q is a directory", o.Dst)
		}
		err = os.Remove(o.Dst)
		if err != nil {
			return fmt.Errorf("RemoteOptions.Exec remove %q: %w", o.Dst, err)
		}
	}

	buf := bytes.NewBufferString(`{{ define "fontcss" }}`)
	buf.Write(b)
	buf.WriteString(`{{ end }}`)

	err = ioutil.WriteFile(o.Dst, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("RemoteOptions.Exec write file %q: %w", o.Dst, err)
	}
	return nil
}
