package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

var ()

type Options struct {
	t      string
	T      *template.Template
	Blog   *BlogOptions
	Mod    *ModOptions
	Img    *ImgOptions
	Remote *RemoteOptions

	// RSS *RSSOptions
	// Map *MapOptions
	// AMP *AMPOptions
	// Webpkg *WebpkgOptions
	// SigXchange *SigXchangeOptions
}

func NewOptions(args []string) (*Options, error) {
	var o Options
	var err error
	f := flag.NewFlagSet("", flag.ExitOnError)
	f.StringVar(&o.t, "t", "templates", "file or flat directory of templates.gohtml")
	f.Parse(args)
	switch f.Arg(1) {
	case "blog":
		o.Blog = NewBlogOptions(f.Args())
	case "mod":
		o.Mod = NewModOptions(f.Args())
	case "img":
		o.Img = NewImgOptions(f.Args())
	case "remote":
		o.Remote = NewRemoteOptions(f.Args())
	default:
		err = fmt.Errorf("NewOptions no known subcommand found: %q", f.Arg(0))
	}
	return &o, err
}
func (o *Options) Exec() error {
	if o.t != "" {
		fi, err := os.Stat(o.t)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("Options.Exec stat %q: %w", o.t, err)
			}
		} else if err == nil {
			if fi.IsDir() {
				o.t = filepath.Join(o.t, "*.gohtml")
			}
			o.T, err = template.ParseFiles(o.t)
			if err != nil {
				return fmt.Errorf("Options.Exec parse file %q: %w", o.t, err)
			}
		}
	}

	switch {
	case o.Blog != nil:
		return o.Blog.Exec(o)
	case o.Mod != nil:
		return o.Mod.Exec(o)
	case o.Img != nil:
		return o.Img.Exec(o)
	case o.Remote != nil:
		return o.Remote.Exec(o)
	}
	return fmt.Errorf("Options.Exec no subcommand to exec")
}

type BlogOptions struct {
	Src string
	Dst string
}

func NewBlogOptions(args []string) *BlogOptions {
	var o BlogOptions
	f := flag.NewFlagSet("blog", flag.ExitOnError)
	f.StringVar(&o.Src, "src", "blog", "source directory")
	f.StringVar(&o.Dst, "dst", "dst", "output directory")
	f.Parse(args)
	return &o
}

func (o *BlogOptions) Exec(opt *Options) error {
	return fmt.Errorf("ErrNotImplemented")
}

type ModData struct {
	Old string
	New string
}

type ModOptions struct {
	Prefix string
	Src    string
	Dst    string
}

func NewModOptions(args []string) *ModOptions {
	var o ModOptions
	f := flag.NewFlagSet("mod", flag.ExitOnError)
	f.StringVar(&o.Prefix, "prefix", "seankhliao.com", "module name prefixes")
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
		mds = append(mds, ModData{string(fields[0]), filepath.Join(o.Prefix, string(fields[1]))})
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
			dst := filepath.Join(o.Dst, md.New+".html")
			err = ioutil.WriteFile(dst, b.Bytes(), 0644)
			if err != nil {
				log.Printf("ModOptions.Exec write %q: %w", dst, err)
			}
		}(md)
	}
	wg.Wait()
	return nil
}

type ImgOptions struct {
	Src string
	Dst string
}

func NewImgOptions(args []string) *ImgOptions {
	var o ImgOptions
	f := flag.NewFlagSet("img", flag.ExitOnError)
	f.StringVar(&o.Src, "src", "img", "source directory")
	f.StringVar(&o.Dst, "dst", "dst", "output directory")
	f.Parse(args)
	return &o
}
func (o *ImgOptions) Exec(opt *Options) error {
	return fmt.Errorf("ErrNotImplemented")
}

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

func main() {
	o, err := NewOptions(os.Args)
	if err != nil {
		log.Print(err)
	}
	err = o.Exec()
	if err != nil {
		log.Print(err)
	}
}
