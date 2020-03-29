---
description: 19th attempt to get a k8s setup working
title: kubernetes cluster 19 part 2
---

### Goals

- ingress
- google sign in protected endpoints

#### _tldr_

config: [seankhliao/kluster @ v0.19.0][kluster]

use with:

- `make create-cluster`
- `make decrypt` // or create appropriate secret files
- `kubectl apply -k .` // maybe repeat a few times if things fail

##### ingress / sign in

- traefik ingress controller
- pomerium auth handler

#### ingress

[traefik][traefik] still easier to setup than ambassador or contour

- `hostPort` because DNS is already pointed at nodes
- consider running as daemonset
- KubernetesCRD as the only provider

set default TLS cert

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: TLSStore
metadata:
  name: default
spec:
  defaultCertificate:
    secretName: api-seankhliao-com-tls
```

#### auth

Running [pomerium][pomerium] because I can't find anything else in the same space?

- _all in one_ mode: does the services mode even work with forward auth?
- maybe config should be in a ConfigMap and inject `idp_*` and `*_secret` through the env?

pomerium config

```yaml
insecure_server: true
grpc_insecure: true
address: :80

authenticate_service_url: https://auth.api.seankhliao.com
forward_auth_url: http://pomerium.networking.svc.cluster.local

idp_provider: "google"
idp_client_id: CHANGE_ME
idp_client_secret: CHANGE_ME

shared_secret: CHANGE_ME
cookie_secret: CHANGE_ME

metrics_address: ":9090"

tracing_provider: jaeger
tracing_debug: true
tracing_jaeger_agent_endpoint: jaeger-agent.monitor.svc.cluster.local:6831

policy:
  - from: https://traefik.api.seankhliao.com
    to: http://example.com
    allowed_users:
      - admin@example.com
```

IngressRoute config

- route both the authenticate service and catch all the `/.pomerium/` paths (css assets)

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: pomerium
spec:
  entryPoints:
    - https
  routes:
    - kind: Rule
      match: Host(`auth.api.seankhliao.com`)
      services:
        - kind: Service
          name: pomerium
          port: 80
    - kind: Rule
      match: PathPrefix(`/.pomerium/`)
      priority: 100
      services:
        - kind: Service
          name: pomerium
          port: 80
  tls: {}
```

Use with Middleware, ex for traefik's dashboard

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: traefik
spec:
  entryPoints:
    - https
  routes:
    - kind: Rule
      match: Host(`traefik.api.seankhliao.com`)
      middlewares:
        - name: auth-traefik
      services:
        - kind: Service
          name: traefik
          port: 9000
  tls: {}
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: auth-traefik
spec:
  forwardAuth:
    address: http://pomerium.networking.svc.cluster.local/?uri=https://traefik.api.seankhliao.com
```

[traefik]: https://docs.traefik.io/
[pomerium]: https://www.pomerium.io/
