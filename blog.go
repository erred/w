package main

import (
	"flag"
	"fmt"
)

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
