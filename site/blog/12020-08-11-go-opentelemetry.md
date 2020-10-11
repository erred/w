---
description: updated opentelemetry stuff
title: go opentelemetry
---

### _OpenTelemetry_

Time passes, APIs change, things break (or not).

go.opentelemetry.io/otel _v0.10.0_

#### _Tracing_

```go
package main

import (
        "context"
        "net/http"

        "go.opentelemetry.io/otel/api/global"
        "go.opentelemetry.io/otel/api/propagation"
        "go.opentelemetry.io/otel/exporters/trace/jaeger"
        sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ExampleTracing() {
        // configure tracer
        // alternatively:
        //     jaegerProvider, flush, err := jaeger.NewExportPipeline(...)
        //     global.SetTraceProvider(jaegerProvider)
        flush, _ := jaeger.InstallNewPipeline(
                jaeger.WithCollectorEndpoint(
                        jaeger.CollectorEndpointFromEnv(),
                ),
                jaeger.WithProcess(jaeger.Process{
                        ServiceName: "my-service",
                }),
                jaeger.WithSDK(&sdktrace.Config{
                        DefaultSampler: sdktrace.ProbabilitySampler(0.1),
                }),
        )
        defer flush()
        tracer := global.Tracer("my-tracer")

        // configure propagation
        // default should be usable?
        // alternatively:
        //     b3Prop := trace.B3{}
        // alternatively:
        //     w3cProp := trace.TraceContext{}
        //     global.SetPropagators(propagation.New(
        //             propagation.WithExtractors(w3cProp),
        //             propagation.WithInjectors(w3cProp),
        //     ))

        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                // extract span from incoming request
                ctx := propagation.ExtractHTTP(r.Context(), global.Propagators(), r.Header)
                ctx, span := tracer.Start(ctx, "span")
                defer span.End()

                // example function
                func(ctx context.Context) {
                        ctx, span := tracer.Start(ctx, "expensive function")
                        defer span.End()
                        // do stuff
                        _ = ctx
                }(ctx)

                // propagate trace into outgoing request
                req, _ := http.NewRequest("GET", "https://opentelemetry.io/", nil)
                propagation.InjectHTTP(ctx, global.Propagators(), req.Header)

                http.DefaultClient.Do(req)
        })
}
```

#### _Metrics_

Limitations:

- only a single set of bucket boundaries are supported
- labels are set at record time, no predeclaring

```go
package main

import (
        "context"
        "net/http"
        "time"

        "go.opentelemetry.io/otel/api/global"
        "go.opentelemetry.io/otel/api/kv"
        "go.opentelemetry.io/otel/api/metric"
        "go.opentelemetry.io/otel/api/unit"
        "go.opentelemetry.io/otel/exporters/metric/prometheus"
)

func ExampleMetric() {
        promExporter, _ := prometheus.InstallNewPipeline(prometheus.Config{
                DefaultHistogramBoundaries: []float64{1, 5, 10, 50, 100},
        })
        http.Handle("/metrics", promExporter)

        meter := global.Meter("service_name")

        counter0 := metric.Must(meter).NewFloat64Counter("counter0",
                metric.WithDescription("hello world"),
                metric.WithUnit(unit.Bytes),
        )
        counter1 := metric.Must(meter).NewFloat64Counter("counter1")

        hist0 := metric.Must(meter).NewFloat64ValueRecorder("hist0")
        hist1 := metric.Must(meter).NewFloat64ValueRecorder("hist1")

        go func() {
                var i int
                for range time.NewTicker(500 * time.Millisecond).C {
                        i++
                        ctx := context.Background()
                        counter0.Add(ctx, float64(i))
                        counter1.Bind(kv.String("key1", "value1"), kv.Int("key2", i%10)).Add(ctx, float64(i))
                        hist0.Record(ctx, float64(i))
                        hist1.Record(ctx, float64(i), kv.String("key1", "value1"), kv.Int("key2", i%10))

                        // alternatively
                        //     meter.RecordBatch(
                        //         ctx,
                        //         []kv.KeyValue{kv.String("key1", "value1"), kv.Int("key2", i%10)},
                        //         counter1.Measurement(float64(i)),
                        //         hist1.Measurement(float64(i)),
                        //     )
                }
        }()

        http.ListenAndServe(":8080", nil)
}
```
