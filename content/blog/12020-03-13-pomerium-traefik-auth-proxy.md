---
description: using pomerium as traefik forward auth
title: pomerium traefik auth proxy
---

### Goals

- docker deployment
- traefk load balancing
- [pomerium](https://pomerium.io) for auth with google

##### _Extra_ Goals

- prometheus monitoring
- jaeger tracing

#### config

so much configutation

##### docker

needs non-default network for _service discovery_ to work

```bash
docker network create br0
```

##### certificates

so this is run manually,
could probably use something else for automation

certs saved to `/var/certs`, requires account wide cloudflare api keys

```bash
docker run --rm -it --name lego \
    -e CLOUDFLARE_EMAIL=me@example.com \
    -e CLOUDFLARE_API_KEY=CHANGE_ME \
    -v /var/certs:/var/certs \
    goacme/lego \
    --dns cloudflare \
    --path /var/certs \
    --email 'me@example.com' \
    --domains 'example.com' \
    --domains '*.example.com' \
    run
```

##### traefik

_tls_ termination and service discovery, expose ports `80` and `443`

```bash
docker run -d --rm --name traefik \
    --network br0 \
    -p 80:80 \
    -p 443:443 \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /var/certs/certificates:/var/certs/certificates \
    -v /var/traefik:/etc/traefik \
    -l 'traefik.enable=true' \
    -l 'traefik.http.routers.https-redirect.entrypoints=http' \
    -l 'traefik.http.routers.https-redirect.rule=HostRegexp(`{any:.*}`)' \
    -l 'traefik.http.routers.https-redirect.middlewares=https-redirect' \
    -l 'traefik.http.middlewares.https-redirect.redirectscheme.scheme=https' \
    traefik:v2.2
```

`traefik.conf`: could be replaced by cli flags

```yaml
global:
  checkNewVersion: false
  sendAnonymousUsage: false

entrypoints:
  traefik:
    address: ":9000"
  http:
    address: ":80"
  https:
    address: ":443"
    forwardedHeaders:
      insecure: true

providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
  file:
    filename: "/etc/traefik/dynamic.yaml"

metrics:
  prometheus: {}

ping: {}

api:
  insecure: true
  dashboard: true

log:
  format: "json"
  level: "INFO"

accessLog:
  format: "json"

tracing:
  jaeger:
    samplingServerURL: "http://jaeger:5778/sampling"
    localAgentHostPort: "jaeger:6831"
```

`dynamic.conf`: no other way of configuring certs

```yaml
tls:
  certificates:
    - certFile: /var/certs/certificates/example.com.crt
      keyFile: /var/certs/certificates/example.com.key
  options:
    default:
      minVersion: VersionTLS13
```

##### pomerium

forward auth provider for traefik

- no service account necessary for personal accounts
- `policy.to` is not used in forward auth mode
- `policy.set_request_headers` is not set
- email headers aren't set so grafana auth.proxy won't work
  _workaround by setting static headers in traefik_
- jaeger tracing is not extracted or inserted to requests
- `authenticate_service_url` + `authenticate_callback_path` together must match
  _API & Services > Credentials > OAuth 2.0 Client IDs > Authorized redirect URIs_
  in the gcloud console

```bash
docker run -d --rm --name pomerium \
    --network br0 \
    -l 'traefik.enable=true' \
    -l 'traefik.http.routers.auth.tls=true' \
    -l 'traefik.http.routers.auth.entrypoints=https' \
    -l 'traefik.http.routers.auth.rule=Host(`auth.example.com`)' \
    -l 'traefik.http.routers.auth.service=auth' \
    -l 'traefik.http.services.auth.loadbalancer.server.port=80' \
    pomerium/pomerium:v0.6.2
```

`config.yaml`: ignore the lacklustre documentation, look at code
[options.go](https://github.com/pomerium/pomerium/blob/master/config/options.go)
and

```yaml
log_level: debug
address: :80
insecure_server: true

authenticate_service_url: https://auth.example.com
authenticate_callback_path: /oauth2/callback

shared_secret: CHANGE_ME
cookie_secret: CHANGE_ME

metrics_address: ":9090"

tracing_provider: jaeger
tracing_debug: true
tracing_jaeger_agent_endpoint: jaeger:6831

forward_auth_url: http://pomerium:443

idp_provider: "google"
idp_client_id: CHANGE_ME
idp_client_secret: CHANGE_ME

policy:
  - from: https://jaeger.example.com
    allowed_users:
      - me@example.com
```

##### httpbin

example unprotected endpoint

```bash
docker run -d --rm --name httpbin \
    --network br0 \
    -l 'traefik.enable=true' \
    -l 'traefik.http.routers.httpbin.tls=true' \
    -l 'traefik.http.routers.httpbin.entrypoints=https' \
    -l 'traefik.http.routers.httpbin.rule=Host(`httpbin.example.com`)' \
    -l 'traefik.http.routers.httpbin.service=httpbin' \
    -l 'traefik.http.services.httpbin.loadbalancer.server.port=80' \
    kennethreitz/httpbin
```

##### jaeger

trace all the things,
example protected endpoint

```bash
docker run -d --rm --name jaeger \
    --network br0 \
    -e SPAN_STORAGE_TYPE=badger \
    -e BADGER_EPHEMERAL=false \
    -e BADGER_DIRECTORY_VALUE=/var/jaeger/data \
    -e BADGER_DIRECTORY_KEY=/var/jaeger/key \
    -v /var/jaeger:/var/jaeger \
    -l 'traefik.enable=true' \
    -l 'traefik.http.routers.jaeger.tls=true' \
    -l 'traefik.http.routers.jaeger.entrypoints=https' \
    -l 'traefik.http.routers.jaeger.rule=Host(`jaeger.example.com`)' \
    -l 'traefik.http.routers.jaeger.middlewares=jaeger' \
    -l 'traefik.http.routers.jaeger.service=jaeger' \
    -l 'traefik.http.middlewares.jaeger.forwardauth.trustForwardHeader=true' \
    -l 'traefik.http.middlewares.jaeger.forwardauth.authResponseHeaders=x-pomerium-jwt-assertion' \
    -l 'traefik.http.middlewares.jaeger.forwardauth.address=http://pomerium:443/?uri=https://jaeger.example.com' \
    -l 'traefik.http.services.jaeger.loadbalancer.server.port=16686' \
    jaegertracing/all-in-one:1.17
```

