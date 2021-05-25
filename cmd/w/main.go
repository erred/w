package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.seankhliao.com/w/v15/internal/static"
	"go.seankhliao.com/w/v15/internal/webserver"
	"k8s.io/klog/v2/klogr"
)

func main() {
	var o webserver.Options
	o.InitFlags(flag.CommandLine)
	flag.Parse()

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	l := klogr.New()

	hmux, err := newHttp(l, static.S)
	if err != nil {
		l.Error(err, "setup http")
		os.Exit(1)
	}

	o.Logger = l
	o.Handler = corsAllowAll(hmux)

	webserver.New(ctx, &o).Run(ctx)
}
