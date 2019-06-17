title = micro boilerplate
date = 2019-06-19
desc = leveraging the kubernetes environment

---

k8s is a platform,
with _features_ that can be used to make things run more smoothly

# lifecycle signals

`SIGINT` and `SIGKILL` are the 2 main ones to care about,
pod termination:

1. `preStop` hook is called
2. `SIGINT` is called
3. after `terminationGracePeriodSeconds`
4. `SIGKILL` is called

```go
func main() {
        // register handlers

        // spin off server
        svr := &http.Server{}
        go svr.ListenAndServe()

        // block on waiting for signal
        sigs := make(chan os.Signal)
        signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL)
        <-sigs

        // handle graceful shutdown
        svr.Shutdown(context.Background())
}
```

# healthcheck

_liveliness_: you can only die once **(unlike cats)**,
pod is restarted if it fails `livelinessProbe`

_readiness_: transient service unavailability should be signalled here,
failure will mean _Service_ won't route to the pod,
sometimes may be easier to just fail a liveliness and restart

`deployment.yaml`:

```yaml
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      ...
      containers:
        - name: liveliness
          ...
          livelinessProbe:
            httpGet:
              path: /health     # /healthz
              port: http
              initialDelaySeconds: 3
              periodSeconds: 3
          readinessProbe:
            httpGet:
              path: /ready      # /readyz
              port: http
              initialDelaySeconds: 3
              periodSeconds: 3
```

# logging

k8s reads structured logs,
[fluentd](https://github.com/fluent/fluentd) is common

there is a tradeoff between structure logging for machines and text logs for humans

with [logrus](https://github.com/sirupsen/logrus)

```go
import (
        log "github.com/sirupsen/logrus"
)

func init() {
        switch os.Getenv("LOG_LEVEL") {
        case "DEBUG":
                log.SetLevel(log.DebugLevel)
        case "INFO":
                log.SetLevel(log.InfoLevel)
        case "ERROR":
                fallthrough
        default:
                log.SetLevel(log.ErrorLevel)
        }

        switch os.Getenv("LOG_FORMAT") {
        case "JSON":
                log.SetFormatter(log.JSONFormatter)
        default:
                log.SetFormatter(log.TextFormatter)
        }
}

func main() {
        log.Debugf("...")
        log.Printf("...") // equiv to log.Infof
        log.Errorf("...")
}
```

# metrics

metrics endpoint works with [prometheus](https://github.com/prometheus/prometheus)

prometheus has a [go client library](https://github.com/prometheus/client_golang)

```go
import (
        "github.com/prometheus/client_golang/prometheus/promhttp"
)
func main() {
        // default metrics
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)
}
```

# tracing

see [opentracing](https://opentracing.io/), also [jaeger](https://github.com/jaegertracing/jaeger)

basically `tracer := opentracing.Tracer`

and `span := tracer.StartSpan()` and `span.Finish()`

## additional

[go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)
