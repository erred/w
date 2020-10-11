---
description: 19th attempt to get a k8s setup working
title: kubernetes cluster 19 part 3
---

### Goals

- monitoring

#### _tldr_

config: [seankhliao/kluster @ v0.19.0][kluster]

use with:

- `make create-cluster`
- `make decrypt` // or create appropriate secret files
- `kubectl apply -k .` // maybe repeat a few times if things fail

##### monitoring

- prometheus metrics
- promtail + loki logs
- grafana visualization
- jaeger tracing

#### prometheus

[prometheus][prometheus] works like magic if you copy the giant kubernetes scrape config from somewhere

- node-expoter: stats about the k8s node (host machine)
- kube-state-metrics: stats about k8s runtime

annotations to specify what to scrape

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
spec:
  template:
    metadata:
      labels:
        app: example
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9000"
        prometheus.io/path: "/metrics"
    spec: ...
```

#### promtail

[promtail][promtail] gets all the logs from all the pods, also copy config block from somewhere

#### loki

[loki][loki] collects all the logs from promtail, who knows why this needs to be a separate service

#### grafana

[grafana][grafana] is where all the data ends up as charts

- grafana has its own user system
- TODO: create charts

grafana config

- trust header to create user
- provisioning

```ini
[security]
disable_initial_admin_creation = true

[users]
allow_sign_up = false
auto_assign_org = true
auto_assign_org_role = Admin

[auth.proxy]
enabled = true
header_name = X-User-Email
header_property = email
auto_sign_up = true

[analytics]
check_for_updates = false

[log]
mode = console
[log.console]
format = json

[paths]
data = /var/lib/grafana/data
logs = /var/log/grafana
plugins = /var/lib/grafana/plugins
provisioning = /etc/grafana/provisioning

[tracing.jaeger]
address = jaeger-agent:6831
```

extra routing and middleware because pomerium can't do it

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: grafana
spec:
  entryPoints:
    - https
  routes:
    - kind: Rule
      match: Host(`grafana.api.seankhliao.com`)
      middlewares:
        - name: auth-grafana
        - name: auth-grafana-email
      services:
        - kind: Service
          name: grafana
          namespace: monitor
          port: 80
  tls: {}
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: auth-grafana
spec:
  forwardAuth:
    address: http://pomerium.networking.svc.cluster.local/?uri=https://grafana.api.seankhliao.com
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: auth-grafana-email
spec:
  headers:
    customRequestHeaders:
      X-User-Email: admin@api.seankhliao.com
```

#### jaeger

[jaeger][jaeger] 1 of 2 competing tracing standards

- other is Zipkin
- OpenTelemetry merges OpenTracing and OpenCensus
- run in _all in one_ mode: who has time for a database
- app -> jaeger-agent -> jaeger-collector -> jaeger-query

traefik

```yaml
tracing:
  jaeger:
    samplingServerURL: "http://jaeger-agent.monitor.svc.cluster.local:5778/sampling"
    localAgentHostPort: "jaeger-agent.monitor.svc.cluster.local:6831"
    gen128Bit: true
```

grafana

```ini
[tracing.jaeger]
address = jaeger-agent:6831
```

pomerium

- only works internally (doesn't coordinate with other services)
- has the most useless name (all) in all in one mode

```yaml
tracing_provider: jaeger
tracing_debug: true
tracing_jaeger_agent_endpoint: jaeger-agent.monitor.svc.cluster.local:6831
```

[jaeger]: https://www.jaegertracing.io/
[grafana]: https://grafana.com/
[loki]: https://github.com/grafana/loki/tree/master/docs
[promtail]: https://github.com/grafana/loki/blob/master/docs/clients/promtail/configuration.md
[prometheus]: https://prometheus.io/
[kluster]: https://github.com/seankhliao/kluster/tree/v0.19.1
