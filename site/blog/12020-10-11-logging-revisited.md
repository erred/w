---
description: another look at logging
title: logging revisited
---

### _logging_ revisited

last time [i looked](/blog/12020-04-12-go-structured-logging/)
i chose [rs/zerolog](https://github.com/rs/zerolog)
as my preferred logging library.

Having used it for 6 months,
I have complaints :)

#### _complaints_

Inefficient / unordered logfmt:
zerolog efficiently creates json,
but unmarshals it into a map for logfmt,
losing order and efficiency.

Timestamp is at the end:
Creating a logger with timestamp adds the timestamp as the last field,
understandable since zerolog needs to pass along
partially created json for a potentially long time,
but still annoying to read as a human.

Long log chains:
I liked this api, it was strictly typed,
but logging even a few things result in very long lines with low density.

Interop with stdlib logger:
Some things `net/http.Server` take a standard `log.Logger`
and it's somewhat not straightforward to create one from zerolog.

#### _criteria_

choosing things again:

_levelled_:
I realise I basically only use `INFO` and `ERROR`.
Ocasionally i'll use `DEBUG` or `TRACE` interchangeably.
Maybe if combined with
[Dave Cheney's blog post](https://dave.cheney.net/2015/11/05/lets-talk-about-logging)
I'll someday be convinced to just do levelled logging (v=0,1,2),
but I still like seeing the distinction
between purely informational and error handling lines.

_structured_:
I still want both json and logfmt output

_speed_:
less concerned about efficiency,
more about getting the output I want

_color_:
don't really care for this to be built in,
it's going to get disabled when run in a remote system anyway,
better to get a local colorizer

#### _list_

Selected from [awesome-go#logging](https://github.com/avelino/awesome-go#logging)

[hashicorp/logutils](https://pkg.go.dev/github.com/hashicorp/logutils)
looks interesting, basically just adding filtering to stdlib log,
less efficient in that you actually have to create the entire log line first.

```go
lf := &logutils.LevelFilter{
        Levels: []logutils.Level{"DEBUG", "INFO", "ERROR"},
        MinLevel: "INFO",
        Write: os.Stderr,
}

log.SetOutput(lf)
log.Printf("[INFO] oops err=%v", err)

// alternative
debug := log.New(lf, "[DEBUG]", log.Ldate|log.Ltime|log.Lmsgprefix)
info := log.New(lf, "[INFO]", log.Ldate|log.Ltime|log.Lmsgprefix)

debug.Printf("oops err=%v", err)
info.Printf("hmm val=%v", 1)
```

[k8s.io/klog/v2](https://pkg.go.dev/k8s.io/klog/v2)
klog v2 gains structured logging with the `*S` variants

```go
klog.InfoS("hello world", "foo", "bar", "x", 1)
Info.ErrorS(err, "oops", "foo", "bar", "y", 2)
```

output, still need to get used to `Lmmdd hh:mm:ss.uuuuuu threadid file:line]`,
especially the unbalanced `]` and the `threadid` is useless in go

```txt
I1012 20:32:36.250791   24371 main.go:11] "hello world" foo="bar" x=1
E1012 20:32:36.250878   24371 main.go:12] "oops" err="an error" foo="bar" y=2
```

#### _idea_

my basic implementation: [go.seankhliao.com/slog](https://pkg.go.dev/go.seankhliao.com/slog)

```go
l := slog.NewText(os.Stderr)
l.Info("hello", "foo", "bar")
l.Error(errors.New("an error"), "oops", "hello", "world")

l = slog.NewJSON(os.Stderr)
l.Info("hello", "foo", "bar")
l.Error(errors.New("an error"), "oops", "hello", "world")
```

output

```txt
2020-10-12T22:05:04+02:00 INF msg="hello" foo="bar"
2020-10-12T22:05:04+02:00 ERR msg="oops" err="an error" hello="world"
{"time":"2020-10-12T22:05:04+02:00", "level":"INF", "msg":"hello", "foo":"bar"}
{"time":"2020-10-12T22:05:04+02:00", "level":"ERR", "msg":"oops", "err":"an error", "hello":"world"}
```
