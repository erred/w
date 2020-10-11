---
description: running TLS to the backend with traefik
title: traefik backend tls
---

### _TLS_ to the backend

Situation: you have traefik as your reverse proxy / API gateway,
traefik terminates TLS for client connections,
but you also want TLS for the traefik-backend part.

_tldr_: possible securely with static backends, insecurely with dynamic (eg k8s)

#### _traefik_

When setting the `scheme` to `https` for backend services,
traefik uses TLS for connecting directly
to the IP address without `ServerName` in the `ClientHello`.
This causes issues since your certificate now needs to be signed for the
IP address (ephemeral on kubernetes) instead of a stable domain name
(any of `$service`, `$service.$namespace`, `$service.$namespace.svc`,
`$service.$namespace.svc.cluster.local` would have been fine).

This appears unlikely to be fully fixed soon,
[partial fix in v2.4](https://github.com/traefik/traefik/issues/4835).

#### _setup_

Using internal Certificate Authority, generate CA cert+key (ex using mkcert).
Inject CA cert into traefik container and run traefik with static config:

```yaml
serversTransport:
  rootCAs:
    - /etc/internal-ca/ca.crt
```

Generate backend cert+key signed by CA key for the IP address and run server with it.
Ex with Cert Manager:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: http-server-internal
spec:
  secretName: http-server-internal
  duration: 2160h
  renewBefore: 360h
  ipAddresses:
    # testing only
    # guessing whick IP the pod will have
    - "10.205.0.44"
    - "10.205.0.45"
    - "10.205.0.46"
    - "10.205.0.47"
    - "10.205.0.48"
    - "10.205.0.49"
    - "10.205.0.50"
  issuerRef:
    name: internal-ca
    kind: ClusterIssuer
```

Setup and IngressRoute to point traefik to the service

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: http-server
spec:
  entryPoints:
    - https
  routes:
    - kind: Rule
      match: Host(`http-server.example.com`)
      services:
        - kind: Service
          name: http-server
          port: 443
          scheme: https
  tls: {}
```
