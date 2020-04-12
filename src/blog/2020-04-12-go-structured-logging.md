---
description: comparing structured logging libraries in go
title: go structured logging
---

### _structured_ logging

who wants to read a mess of words anyways

#### tldr

use `github.com/rs/zerolog` because I like the api and the output format

#### about

##### criteria

- structured/levelled
- json output
- logfmt/txt output

#### Benchmark

code [gist](https://gist.github.com/seankhliao/259ab478721d84ad6f2f5935f7fe052f)

```txt
» go test -bench . -benchmem
goos: linux
goarch: amd64
pkg: t
BenchmarkApexTxt-4                     305184              3541 ns/op            1228 B/op              20 allocs/op
BenchmarkApexJSON-4                    275408              4287 ns/op            1568 B/op              23 allocs/op
BenchmarkGokitTxt-4                   1000000              1154 ns/op             468 B/op               5 allocs/op
BenchmarkGokitJSON-4                   319072              3619 ns/op            1232 B/op              22 allocs/op
BenchmarkLog15JSON-4                   172742              7028 ns/op            2208 B/op              33 allocs/op
BenchmarkLog15Txt-4                    324510              3578 ns/op            1376 B/op              20 allocs/op
BenchmarkLogrusJSON-4                  203190              5791 ns/op            1974 B/op              35 allocs/op
BenchmarkLogrusTxt-4                   260695              4694 ns/op            1056 B/op              23 allocs/op
BenchmarkZapJSON-4                    1000000              1057 ns/op             256 B/op               1 allocs/op
BenchmarkZapTxt-4                      724332              1560 ns/op             368 B/op               5 allocs/op
BenchmarkZapSugarJSON-4                750585              1370 ns/op             512 B/op               1 allocs/op
BenchmarkZapSugarTxt-4                 605372              1855 ns/op             624 B/op               5 allocs/op
BenchmarkZerologJSON-4                2005058               588 ns/op               0 B/op               0 allocs/op
BenchmarkZerologTxt-4                   98937             11678 ns/op            2873 B/op              93 allocs/op
BenchmarkKlogInterface-4               736012              1631 ns/op             216 B/op               2 allocs/op
BenchmarkKlogFormat-4                  797425              1501 ns/op             216 B/op               2 allocs/op
BenchmarkStdlibInterface-4            1978610               607 ns/op              48 B/op               1 allocs/op
BenchmarkStdlibFormat-4               2119566               565 ns/op              48 B/op               1 allocs/op
PASS
ok          t        23.216s
```

#### loggers

- `github.com/apex/log`
- `github.com/go-kit/kit/log`
- `github.com/inconshreveable/log15`
- `github.com/rs/zerolog`
- `github.com/sirupsen/logrus`
- `go.uber.org/zap`

not structured:

- `k8s.io/klog` (fork of glog)
- `log` (stdlib)

##### apex

code:

```go
import (
        apex "github.com/apex/log"
        "github.com/apex/log/handlers/json"
        "github.com/apex/log/handlers/logfmt"
)

func Apex() {
        fmt.Println("\nApex")
        l := &apex.Logger{Handler: json.New(os.Stdout), Level: apex.InfoLevel}
        l.WithFields(apex.Fields{"bool": true, "int": 1, "str": "s", "err": err}).Info("apex json")

        l = &apex.Logger{Handler: logfmt.New(os.Stdout), Level: apex.InfoLevel}
        l.WithFields(apex.Fields{"bool": true, "int": 1, "str": "s", "err": err}).Info("apex txt")
}
```

output:

```txt
Apex
{"fields":{"bool":true,"err":{},"int":1,"str":"s"},"level":"info","timestamp":"2020-04-12T22:58:32.105258571+02:00","message":"apex json"}
timestamp=2020-04-12T22:58:32.105372627+02:00 level=info message="apex txt" bool=true err="an err" int=1 str=s
```

##### Gokit

code:

```go
import (
        kitlog "github.com/go-kit/kit/log"
)

func Gokit() {
        fmt.Println("\nGokit")
        l := kitlog.NewJSONLogger(os.Stdout)
        l.Log("gokit json", "bool", true, "int", 1, "str", "s", "err", err)

        l = kitlog.NewLogfmtLogger(os.Stdout)
        l.Log("gokit txt", "bool", true, "int", 1, "str", "s", "err", err)
}
```

output:

```txt
Gokit
{"1":"str","an err":"(MISSING)","gokit json":"bool","s":"err","true":"int"}
gokittxt=bool true=int 1=str s=err
```

##### Log15

code:

```go
import (
        "github.com/inconshreveable/log15"
)

func Log15() {
        fmt.Println("\nLog15")
        l := log15.New()
        l.SetHandler(log15.StreamHandler(os.Stdout, log15.JsonFormat()))
        l.Info("log15 json", "bool", true, "int", 1, "str", "s", "err", err)

        l = log15.New()
        l.SetHandler(log15.StreamHandler(os.Stdout, log15.LogfmtFormat()))
        l.Info("log15 json", "bool", true, "int", 1, "str", "s", "err", err)
}
```

output:

```txt
Log15
{"bool":true,"err":"an err","int":1,"lvl":"info","msg":"log15 json","str":"s","t":"2020-04-12T22:58:32.105473094+02:00"}
t=2020-04-12T22:58:32+0200 lvl=info msg="log15 json" bool=true int=1 str=s err="an err"
```

##### Logrus

code:

```go
import (
        "github.com/sirupsen/logrus"
)

func Logrus() {
        fmt.Println("\nLogrus")
        l := &logrus.Logger{Out: os.Stdout, Formatter: &logrus.JSONFormatter{}, Level: logrus.InfoLevel}
        l.WithFields(logrus.Fields{"bool": true, "int": 1, "str": "s", "err": err}).Info("logrus json")

        l = &logrus.Logger{Out: os.Stdout, Formatter: &logrus.TextFormatter{}, Level: logrus.InfoLevel}
        l.WithFields(logrus.Fields{"bool": true, "int": 1, "str": "s", "err": err}).Info("logrus txt")
}
```

output:

```txt
Logrus
{"bool":true,"err":"an err","int":1,"level":"info","msg":"logrus json","str":"s","time":"2020-04-12T22:58:32+02:00"}
INFO[0000] logrus txt                                    bool=true err="an err" int=1 str=s
```

##### Zap

code:

```go
import(
        "go.uber.org/zap"
)

func Zap() {
        fmt.Println("\nZap")
        l, _ := zap.NewProduction()
        l.Info("zap json", zap.Bool("bool", true), zap.Int("int", 1), zap.String("str", "s"), zap.Error(err))

        l, _ = zap.NewDevelopment()
        l.Info("zap txt", zap.Bool("bool", true), zap.Int("int", 1), zap.String("str", "s"), zap.Error(err))
}
```

output:

```txt
Zap
{"level":"info","ts":1586725112.1056206,"caller":"testrepo-82/main.go:78","msg":"zap json","bool":true,"int":1,"str":"s","error":"an err"}
2020-04-12T22:58:32.105+0200        INFO        testrepo-82/main.go:81        zap txt        {"bool": true, "int": 1, "str": "s", "error": "an err"}
```

##### Zap Sugared

code:

```go
import(
        "go.uber.org/zap"
)

func ZapSugar() {
        fmt.Println("\nZap Sugar")
        l, _ := zap.NewProduction()
        l.Sugar().Infow("zap sugar json", "bool", true, "int", 1, "str", "s", "err", err)

        l, _ = zap.NewDevelopment()
        l.Sugar().Infow("zap sugar txt", "bool", true, "int", 1, "str", "s", "err", err)
}
```

output:

```txt
Zap Sugar
{"level":"info","ts":1586725112.1056879,"caller":"testrepo-82/main.go:87","msg":"zap sugar json","bool":true,"int":1,"str":"s","err":"an err"}
2020-04-12T22:58:32.105+0200        INFO        testrepo-82/main.go:90        zap sugar txt        {"bool": true, "int": 1, "str": "s", "err": "an err"}
```

##### Zerolog

code:

```go
import(
        "github.com/rs/zerolog"
)

func Zerolog() {
        fmt.Println("\nZerolog")
        l := zerolog.New(os.Stdout).With().Timestamp().Logger()
        l.Info().Bool("bool", true).Int("int", 0).Str("str", "s").Err(err).Msg("zerolog json")

        l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
        l.Info().Bool("bool", true).Int("int", 1).Str("str", "s").Err(err).Msg("zerolog txt")
}
```

output:

```txt
Zerolog
{"level":"info","bool":true,"int":0,"str":"s","error":"an err","time":"2020-04-12T22:58:32+02:00","message":"zerolog json"}
2020-04-12T22:58:32+02:00 INF zerolog txt error="an err" bool=true int=1 str=s
```

##### Klog

code:

```go
import(
        "k8s.io/klog"
)

func Klog() {
        fmt.Println("\nKlog")
        klog.Info("klog", "bool", true, "int", 1, "str", "s", "err", err)
        klog.Infof("klog bool=%v int=%v str=%v err=%v\n", true, 1, "s", err)
}
```

```txt
Klog
I0412 22:58:32.105824  258002 main.go:104] klogbooltrueint1strserran err
I0412 22:58:32.105840  258002 main.go:105] klog bool=true int=1 str=s err=an err
```

##### Log

code:

```go
import(
        "log"
)

func Std() {
        fmt.Println("\nstdlib")
        log.Println("stdlib", "bool", true, "int", 1, "str", "s", "err", err)
        log.Printf("klog bool=%v int=%v str=%v err=%v\n", true, 1, "s", err)
}
```

output:

```txt
stdlib
2020/04/12 23:01:38 stdlib bool true int 1 str s err an err
2020/04/12 23:01:38 stdlib bool=true int=1 str=s err=an err
```
