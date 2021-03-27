---
title: ingress-nginx opentracing
description: docs are few and far between
---

### _ingress-nginx_

Don't ask me why I'm using the default [ingress-nginx](https://kubernetes.github.io/ingress-nginx/),
I just have to.
Anyway, I want to try out its integration with [jaeger](https://www.jaegertracing.io/)

#### _config_

##### _minimal_

Annoyingly I spent 4 hours only to realize the reason my ingress was crashing
was because the docs show the latest master but I'm running the latest release,
and the `jaeger-endpoint` config option was very recent.

```yaml
data:
  enable-opentracing: "true"
  # agent host (udp/6831) adjust accordingly
  jaeger-collector-host: jaeger.jaeger.svc.cluster.local

  # takes precedence over collector-host
  # only available after 0.44.0
  # jaeger-endpoint: "http://jaeger.jaeger.svc.cluster.local:14268/api/traces"
```

##### _other_

For reasons, the default config is a bit useless: each request generates 2 spans:
both with the name service name **nginx** and the operation name being the Ingress path,
which for most people is just **/**.

Thankfully, they can be customized, especially useful are the available [log fields](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/log-format/)

```yaml
data:
  # rename the service
  jaeger-service-name: nginx-default

  # outer span
  opentracing-operation-name: "$request_method $host"
  # inner span
  opentracing-location-operation-name: "$namespace $service_name"

  # don't kill jaeger
  jaeger-sampler-type: ratelimiting
  jaeger-sampler-param: "5"
```
