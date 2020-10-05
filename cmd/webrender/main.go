package main

import (
	"flag"
	"fmt"
	"os"

	"go.seankhliao.com/com-seankhliao/v12/render"
	"go.seankhliao.com/com-seankhliao/v12/render/style"
)

func main() {
	var o render.Options
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	o.InitFlags(fs)
	fs.Parse(os.Args[1:])

	o.Template = style.Template

	err := render.Process(o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render: %v", err)
		os.Exit(1)
	}
}
