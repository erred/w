---
title: klog output
description: writing out logs has never been so hard...
---

### _klog_

[klog](https://pkg.go.dev/k8s.io/klog/v2) is interesting:
it implements
[logr.Logger](https://pkg.go.dev/github.com/go-logr/logr#Logger)
with [klogr](https://pkg.go.dev/k8s.io/klog/klogr)
but also [accepts](https://pkg.go.dev/k8s.io/klog/v2#SetLogger)
a `logr.Logger` as the output.

Anyway, here's some commented code to show how the setups interact.
Passing `klogr->klog->some_other_logr` looks inefficient and doesn't get you the best results anyways...

```go
package main

import (
	"errors"

	"github.com/butonic/zerologr"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
)

func main() {
	zl := zerologr.New() // an implementation of logr.Logger
	err := errors.New("some: error")

  // global default klog
	klog.Errorln("hello world", "foo", "bar")
	// E0617 23:38:11.432545   83114 main.go:16] hello world foo bar
	klog.ErrorS(err, "oopsie daisy", "foo", "bar")
	// E0617 23:38:11.432664   83114 main.go:17] "oopsie daisy" err="some: error" foo="bar"


  // klogr with default serialization (pass single string to klog)
	kl := klogr.NewWithOptions(klogr.WithFormat(klogr.FormatSerialize))
	kl = kl.WithName("app").WithName("comp1")
	kl.Info("hello world", "foo", "bar")
	// I0617 23:38:11.432704   83114 main.go:21] app/comp1 "msg"="hello world"  "foo"="bar"
	kl.Error(err, "oopsie daisy", "foo", "bar")
	// E0617 23:38:11.432720   83114 main.go:22] app/comp1 "msg"="oopsie daisy" "error"="some: error"  "foo"="bar"


  // klogr passing serialization to klog
	kl = klogr.NewWithOptions(klogr.WithFormat(klogr.FormatKlog))
	kl = kl.WithName("app").WithName("comp1")
	kl.Info("hello world", "foo", "bar")
	// I0617 23:38:11.432736   83114 main.go:26] "app/comp1: hello world" foo="bar"
	kl.Error(err, "oopsie daisy", "foo", "bar")
	// E0617 23:38:11.432749   83114 main.go:27] "app/comp1: oopsie daisy" err="some: error" foo="bar"


  // use a logr.Logger for klog output
  // default output
	klog.SetLogger(zl)
	klog.Errorln("hello world", "foo", "bar")
	// {"level":"error","time":"2021-06-17T23:38:11+02:00","message":"hello world foo bar\n"}
	klog.ErrorS(err, "oopsie daisy", "foo", "bar")
	// {"level":"error","error":"some: error","foo":"bar","time":"2021-06-17T23:38:11+02:00","message":"oopsie daisy"}


  // klogr with default serialization (pass single string to klog)
	kl = klogr.NewWithOptions(klogr.WithFormat(klogr.FormatSerialize))
	kl = kl.WithName("app").WithName("comp1")
	kl.Info("hello world", "foo", "bar")
	// {"level":"info","verbosity":0,"time":"2021-06-17T23:38:11+02:00","message":"app/comp1 \"msg\"=\"hello world\"  \"foo\"=\"bar\"\n"}
	kl.Error(err, "oopsie daisy", "foo", "bar")
	// {"level":"error","time":"2021-06-17T23:38:11+02:00","message":"app/comp1 \"msg\"=\"oopsie daisy\" \"error\"=\"some: error\"  \"foo\"=\"bar\"\n"}


  // klogr passing serialization to klog
	kl = klogr.NewWithOptions(klogr.WithFormat(klogr.FormatKlog))
	kl = kl.WithName("app").WithName("comp1")
	kl.Info("hello world", "foo", "bar")
	// {"level":"info","verbosity":0,"foo":"bar","time":"2021-06-17T23:38:11+02:00","message":"app/comp1: hello world"}
	kl.Error(err, "oopsie daisy", "foo", "bar")
	// {"level":"error","error":"some: error","foo":"bar","time":"2021-06-17T23:38:11+02:00","message":"app/comp1: oopsie daisy"}


	// what you would get if you used the logr.Logger directly
	zl = zl.WithName("app").WithName("comp1")
	zl.Info("hello world", "foo", "bar")
	// {"level":"info","verbosity":0,"name":"app/comp1","foo":"bar","time":"2021-06-17T23:52:34+02:00","message":"hello world"}
	zl.Error(err, "oopsie daisy", "foo", "bar")
	// {"level":"error","error":"some: error","name":"app/comp1","foo":"bar","time":"2021-06-17T23:52:34+02:00","message":"oopsie daisy"}
}
```
