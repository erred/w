package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	host   string             // hostname
	t      string             // template file or dir (flat)
	T      *template.Template // parsed templates
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
