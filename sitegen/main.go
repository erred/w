package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

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

// Options contains global options
// such as the template dir
// and host
type Options struct {
	host string             // hostname
	t    string             // template file or dir (flat)
	T    *template.Template // parsed templates

	// creates template
	Remote *RemoteOptions

	// creates content
	Static *StaticOptions
	Blog   *BlogOptions
	Mod    *ModOptions
	Img    *ImgOptions

	// creates metadata
	Sitemap *SitemapOptions

	// RSS *RSSOptions
	// AMP *AMPOptions
	// Webpkg *WebpkgOptions
	// SigXchange *SigXchangeOptions
}

// NewOptions parses global options and decides on subcommand
// remember to register here
func NewOptions(args []string) (*Options, error) {
	var o Options
	var err error
	f := flag.NewFlagSet("", flag.ExitOnError)
	f.StringVar(&o.host, "host", "seankhliao.com", "hostname")
	f.StringVar(&o.t, "t", "templates", "file or flat directory of templates.gohtml")
	f.Parse(args)
	switch f.Arg(1) {
	case "static":
		o.Static = NewStaticOptions(f.Args())
	case "blog":
		o.Blog = NewBlogOptions(f.Args())
	case "mod":
		o.Mod = NewModOptions(f.Args())
	case "img":
		o.Img = NewImgOptions(f.Args())
	case "remote":
		o.Remote = NewRemoteOptions(f.Args())
	case "sitemap":
		o.Sitemap = NewSitemapOptions(f.Args())
	default:
		err = fmt.Errorf("NewOptions no known subcommand found: %q", f.Arg(0))
	}
	return &o, err
}

// Exec parses global templates and executes a subcommand
// remember to register here too
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
			o.T, err = template.ParseGlob(o.t)
			if err != nil {
				return fmt.Errorf("Options.Exec parse file %q: %w", o.t, err)
			}
		}
	}

	switch {
	case o.Static != nil:
		return o.Static.Exec(o)
	case o.Blog != nil:
		return o.Blog.Exec(o)
	case o.Mod != nil:
		return o.Mod.Exec(o)
	case o.Img != nil:
		return o.Img.Exec(o)
	case o.Remote != nil:
		return o.Remote.Exec(o)
	case o.Sitemap != nil:
		return o.Sitemap.Exec(o)
	}
	return fmt.Errorf("Options.Exec no subcommand to exec")
}

func canonicalURL(subpath string) string {
	subpath = strings.TrimSuffix(subpath, ".html")
	subpath = strings.TrimSuffix(subpath, "index")
	subpath = strings.TrimSuffix(subpath, "/")
	if subpath == "" {
		subpath = "/"
	}
	return subpath
}
