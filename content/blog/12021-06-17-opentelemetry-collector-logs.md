---
title: opentelemetry-collector logs
description: collecting logs with otel
---

### _logs_

The formalized version of _printf_ debugging.
Anyway, you run a lot of applications written by a lot of different people.
Their logs are all over the place, what can you do?

Recently the [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/),
with the donation of [stanza / opentelemetry-log-collection](https://github.com/open-telemetry/opentelemetry-log-collection)
gained the ability to [collect logs](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/filelogreceiver),
do some light processing and ship them elsewhere.
It's still early so sort of limited in what it can do.

Let's look at log outputs from some Go logging libraries and how you'd parse them.

#### _std_ log

While the [log](https://pkg.go.dev/log) package does have flags for you to control output,
In the standard config it just adds time.
It also has no structure in the messages:

```
2021/06/17 22:37:44 Something happened foo=bar
```

This can be parsed with:

```yaml
receivers:
  filelog/std:
    include:
      - std.log
    operators:
      - type: regex_parser
        regex: '^(?P<timestamp_field>.{19}) (?P<message>.*)$'
        timestamp:
          parse_from: timestamp_field
          layout_type: strptime
          layout: '%Y/%m/%d %T'
```

Resulting in (otel logging exporter output):

```
Resource labels:
     -> host.name: STRING(eevee)
     -> os.type: STRING(LINUX)
InstrumentationLibraryLogs #0
InstrumentationLibrary
LogRecord #0
Timestamp: 2021-06-17 22:37:44 +0000 UTC
Severity: Undefined
ShortName:
Body: {
     -> message: STRING(Something happened foo=bar)
}
```

#### _json_

There are a lot of structured loggers out there,
I happen to like [zerolog](https://github.com/rs/zerolog).

```
{"level":"error","error":"oops","foo":"bar","time":"2021-06-17T22:38:02+02:00","message":"something bad happened"}
```

json is well supported

_note:_ the
[strptime](https://github.com/observiq/ctimefmt/blob/3e07deba22cf7a753f197ef33892023052f26614/ctimefmt.go#L63)
parser seems to take issue with the timezone for some reason,
so I'm using the [Go time parser](https://pkg.go.dev/time#Parse)

```yaml
receivers:
  filelog/json:
    include:
      - json.log
    include_file_name: false
    operators:
      - type: json_parser
        timestamp:
          parse_from: time
          layout_type: gotime
          layout: 2006-01-02T15:04:05Z07:00
        severity:
          parse_from: level
```

output:

```
Resource labels:
     -> host.name: STRING(eevee)
     -> os.type: STRING(LINUX)
InstrumentationLibraryLogs #0
InstrumentationLibrary
LogRecord #0
Timestamp: 2021-06-17 20:38:02 +0000 UTC
Severity: Error
ShortName:
Body: {
     -> error: STRING(oops)
     -> foo: STRING(bar)
     -> message: STRING(something bad happened)
}
```

#### _klog_ / glog

[klog](https://pkg.go.dev/k8s.io/klog/v2) is kubernetes' standard logger,
mostly based on [glog](https://pkg.go.dev/github.com/golang/glog).
And recent versions have gained support for structured logging.
But this is where we reach the limits of the current log parser.

Unlike [loki](https://github.com/grafana/loki) it doesn't
have a [logfmt](https://brandur.org/logfmt)
[parser](https://grafana.com/docs/loki/latest/logql/#logfmt)
meaning your `key=value` pairs are just stuck there.
All the more reason to use json loggers then...

```
E0617 22:38:02.013247   76356 main.go:57]  "msg"="something bad happened" "error"="oops"
```

config:

```yaml
receivers:
  filelog/klog:
    include:
      - klog.log
    include_file_name: false
    operators:
      - type: regex_parser
        # Lmmdd hh:mm:ss.uuuuuu threadid file:line]
        regex: '^(?P<level>[EI])(?P<timestamp_field>.{20})\s+(?P<threadid>\d+)\s(?P<file>\w+\.go):(?P<line>\d+)]\s+(?P<message>.*)$'
        timestamp:
          parse_from: timestamp_field
          layout: '%m%d %H:%M:%S.%f'
        severity:
          parse_from: level
          mapping:
            error: E
            info: I
```

result

```
Resource labels:
     -> host.name: STRING(eevee)
     -> os.type: STRING(LINUX)
InstrumentationLibraryLogs #0
InstrumentationLibrary
LogRecord #0
Timestamp: 2021-06-17 22:38:02.013247 +0000 UTC
Severity: Error
ShortName:
Body: {
     -> file: STRING(main.go)
     -> line: STRING(57)
     -> msg: STRING("msg"="something bad happened" "error"="oops"  )
     -> threadid: STRING(73779)
}
