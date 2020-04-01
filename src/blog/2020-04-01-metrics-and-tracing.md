---
description: metrics and tracing with prometheus, jaeger, opentracing, opencensus, and opentelemetry
title: metrics and tracing
---

### metrics and tracing

So you created a mess of services that talk to each other
and now you want to know what actually happens.

#### defs

##### metrics are...

counters providing an aggregate view of events

##### tracing is...

records of individual execuation

#### parts

there are 2 parts:

##### collectors

These are seperate applications that run to collect metrics/traces.
Usually the things with name recognition,
ex [prometheus][prom] for metrics and [jaeger][jaeger] for tracing

_TODO_: look at [opentelemetry collector][otc] for collecting all metrics / traces
and reexport them to collectors

##### sdk/client libraries

these run as part of application code and generate/expose metrics/traces.
Subject of ongoing standardization, ex

- [prometheus/client_golang][promgo] for metrics, prometheus official library
- [OpenMetrics][om] for metrics, still not usable
- [OpenCensus][oc] for metrics and tracing
- [jaegertracing/jaeger-client-go][jaegergo] for tracing, jaeger official library
- [OpenTracing][ot] for tracing
- [OpenTelemetry][otel] for tracing and metrics, [merging][merge] OpenTracing and OpenCensus

OpenCensus and OpenTelemetry implement metrics and tracing themselves
and just expose tracing / metrics in a way the collectors can understand,
no official clients required (might not always be true)

#### examples

with different library combinatioons:

- initialize metrics with prometheus
- initialize tracing with jaeger
- export metrics from tracing
- extract trace from incoming request
- inject trace into outgoing request

_note_: there are 3 main methods of trace propagation

- Zipkin B3
- Jaeger Uber
- W3C TraceContext: subject of ongoing standardization

##### Prometheus / Jaeger libraries

using [prometheus/client_golang][promgo] and [jaegertracing/jaeger-client-go][jaegergo]

jaeger officially uses [opentracing/opentracing-go][otgo]

only way to push metrics to prometheus pushgateway

```go
package main

import (
        "net/http"

        "github.com/opentracing/opentracing-go"
        "github.com/opentracing/opentracing-go/ext"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
        "github.com/prometheus/client_golang/prometheus/push"
        "github.com/uber/jaeger-client-go"
        jaegercfg "github.com/uber/jaeger-client-go/config"
        jprom "github.com/uber/jaeger-lib/metrics/prometheus"
)

func ExamplePrometheusJaeger() {
        // stats: prometheus exporter
        http.Handle("/metrics", promhttp.Handler())

        // stats: custom
        counter0 := promauto.NewCounter(prometheus.CounterOpts{
                Name: "myapp_processed_ops_total",
                Help: "The total number of processed events",
        })
        prometheus.MustRegister(counter0)

        // trace: jaeger exporter
        // or FromEnv()
        cfg := jaegercfg.Configuration{
                ServiceName: "service",
                Sampler: &jaegercfg.SamplerConfig{
                        Type:  jaeger.SamplerTypeConst,
                        Param: 1,
                },
        }
        tracer, closer, _ := cfg.NewTracer(jaegercfg.Metrics(jprom.New()))
        defer closer.Close()
        opentracing.InitGlobalTracer(tracer)

        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                // extract span from incoming context
                spanctx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
                span := tracer.StartSpan("child", ext.RPCServerOption(spanctx))
                defer span.Finish()

                // update metrics
                counter0.Inc()

                // propagate trace into outgoing request
                req, _ := http.NewRequest("GET", "https://opencensus.io/", nil)
                tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

                http.DefaultClient.Do(req)
        })

        // serve
        http.ListenAndServe(":8080", nil)

        // or push metrics to pushgateway
        pusher := push.New("push-gateway:1234", "job-name").Gatherer(prometheus.DefaultGatherer)
        pusher.Push()
}
```

#### OpenCensus

using [go.opencensus.io][ocgo]

Metrics looks clunky without increment/decrement.
Also have to replace http client/server

```go
package main

import (
        "net/http"

        "contrib.go.opencensus.io/exporter/jaeger"
        "contrib.go.opencensus.io/exporter/prometheus"
        "go.opencensus.io/plugin/ochttp"
        "go.opencensus.io/stats"
        "go.opencensus.io/stats/view"
        "go.opencensus.io/trace"
)

func ExampleOpenCensus() {
        // stats: prometheus exporter
        promexp, _ := prometheus.NewExporter(prometheus.Options{})
        http.Handle("/metrics", promexp)

        // stats: ochttp defaults
        view.Register(ochttp.DefaultClientViews...)
        view.Register(ochttp.DefaultServerViews...)

        // stats: custom
        float0 := stats.Float64("float0", "A float", "ms")

        // trace: jaeger exporter
        jaegerexp, _ := jaeger.NewExporter(jaeger.Options{
                AgentEndpoint: "localhost:6831",
                ServiceName:   "demo",
        })
        trace.RegisterExporter(jaegerexp)
        trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

        // propagation
        client := &http.Client{Transport: &ochttp.Transport{}}
        handler := &ochttp.Handler{Handler: http.DefaultServeMux}

        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                // extract span from incoming request
                ctx, span := trace.StartSpan(r.Context(), "span-name")
                defer span.End()

                // update metrics
                stats.Record(ctx, float0.M(1.2))

                // propagate trace into outgoing request
                req, _ := http.NewRequest("GET", "https://opencensus.io/", nil)
                req = req.WithContext(ctx)

                client.Do(req)
        })

        http.ListenAndServe(":8080", handler)
}
```

#### OpenTelemetry

using [go.opentelemetry.io/otel][otelgo]

Okayish but api feels unstable and overly complex.
Also doesn't support jaeger native trace propagation(?)

```go
package main

import (
        "net/http"

        "go.opentelemetry.io/otel/api/global"
        "go.opentelemetry.io/otel/api/metric"
        "go.opentelemetry.io/otel/api/propagation"
        "go.opentelemetry.io/otel/exporters/metric/prometheus"
        "go.opentelemetry.io/otel/exporters/trace/jaeger"
        sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ExampleOpenTelemetry() {
        // stats: prometheus exporter
        controller, promhandle, _ := prometheus.InstallNewPipeline(prometheus.Config{})
        defer controller.Stop()
        http.Handle("/metrics", promhandle)

        // stats: custom
        meter := global.Meter("service")
        counter0 := metric.Must(meter).NewInt64Counter("counter0")

        // trace: jaeger exporter
        _, flush, _ := jaeger.NewExportPipeline(
                jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
                jaeger.RegisterAsGlobal(),
                jaeger.WithSDK(&sdktrace.Config{
                        DefaultSampler: sdktrace.AlwaysSample(),
                }),
        )
        defer flush()
        tracer := global.Tracer("service")

        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                // extract span from incoming request
                ctx := propagation.ExtractHTTP(r.Context(), global.Propagators(), r.Header)
                ctx, span := tracer.Start(ctx, "span")
                defer span.End()

                // update metrics
                counter0.Add(ctx, 1)

                // propagate trace into outgoing request
                req, _ := http.NewRequest("GET", "https://opencensus.io/", nil)
                propagation.InjectHTTP(ctx, global.Propagators(), req.Header)

                http.DefaultClient.Do(req)
        })
}
```

[otc]: https://opentelemetry.io/docs/collector/about/
[otgo]: https://github.com/opentracing/opentracing-go
[otelgo]: https://github.com/open-telemetry/opentelemetry-go
[ocgo]: https://github.com/census-instrumentation/opencensus-go
[jaegergo]: https://github.com/jaegertracing/jaeger-client-go
[promgo]: https://github.com/prometheus/client_golang
[merge]: https://medium.com/opentracing/merging-opentracing-and-opencensus-f0fe9c7ca6f0
[ot]: https://opentracing.io/
[oc]: https://opencensus.io/
[otel]: https://opentelemetry.io/
[om]: https://openmetrics.io/
[prom]: https://prometheus.io/
[jaeger]: https://www.jaegertracing.io/
