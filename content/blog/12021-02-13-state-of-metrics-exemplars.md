---
description: what needs to happen to get exemplars in prometheus and grafana
title: state of metrics exemplars
---

### _exemplars_

_tldr:_ not here (yet)

So you can already link between traces and logs,
but what about connecting them to metrics?
That means _exemplars_,
individual data points whith store extra labels,
not just the aggregate data usually see in metrics.

A [public design doc](https://docs.google.com/document/d/1ymZlc9yuTj8GvZyKz1r3KDRrhaOjZ1W1qZVW_5Gj7gA/edit#)
is available for prometheus.

#### _required_ changes

The decision to record exemplars starts at the client,
so, API changes a necessary.
For a counter, it's add, but you store additional labels
(like `traceID`) for this particular point.
This will look something like:

```go
counter.AddWithExemplar(value float64, exemplar prometheus.Labels)
```

Next is to expose the data:
exemplars are part of the [OpenMetrics](https://github.com/OpenObservability/OpenMetrics) spec,
only one exemplar per metric / bucket.

```txt
foo_bucket{le=”0.1”} 8 <timestamp> # {id=abc} 0.043 <timestamp>
```

Since they're comments, they could be ignored,
but what good would that do?
Prometheus needs to store them, so we can query (and graph) them later.
Work is primarily in [prometheus/prometheus#6635](https://github.com/prometheus/prometheus/pull/6635),
which didn't make it in to `v2.25.0`.

Final piece is grafana, which should have basic support in `v7.4.0`,
but it's kind of hard to test when storage doesn't work...
