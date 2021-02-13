---
description: connecting traces and logs in grafana
title: grafana trace to log to trace
---

### _grafana_

The recently released [Grafana](https://github.com/grafana/grafana)
[7.4.0](https://github.com/grafana/grafana/releases/tag/v7.4.0)
includes support for jumping between logs and traces
(_only in the explore tab_).

#### _apps_

your apps need to: include traceids in logs and include labels in spans (can probably be done with otel-collector...)

#### _setup_

Get logs into grafana: promtail to scrape the logs, loki to aggregate / store.
Add Loki as a datasource to grafana,
add a derived field extracting the traceid and linking to jaeger.

Get traces into grafana: jaeger to collect the logs.
Add jaeger as a datasource to grafana,
select loki for trace to logs, add extra tags if you have them
(to make filtering more specific).

#### _play_

Go to the explore tab.

##### _traces_

Currently you can only retrieve a trace if you know the traceID
or you can select from a few in the query selector. No searching (yet).
Trace to log is available as a little button next to the span name.

##### _logs_

Crash course LogQL (see [docs](https://grafana.com/docs/loki/latest/logql/) for details):

```logql
{app="svc"} |= "info" | logfmt | x > 5 | label_format ... | line_format ...
```

- `{app="svc"}`: select logs using labels
- `|=`: initial filter over raw log line (contains, regex)
- `logfmt`: parse into structured data, also supports `json`, `regex`
- `x > 5`: query structured data
- `label_format`: modify labels
- `line_format`: output as raw log line

Link to trace is available once you click open a log line.
A button will show next to to the trace field (configured in the datasource).
