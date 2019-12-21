running traefik v2 on k8s

---

# v2 alpha

[traefik](https://docs.traefik.io/v2.0/) is in its v2 alpha,
so what does it take to get running

## _ingress_ (route)

auto discover services with `--providers.kubernetes` and `--providers.kubernetescrd`,
to discover them from Ingresses and IngressRoutes

## acme

traefik can [automatically get](https://docs.traefik.io/v2.0/https-tls/acme/)
tls certificates from providers such as [let's encrypt](https://letsencrypt.org/),
but they have rate limiting.
also traefik can only store the certs in a file (no shared use)
or a key-value store

instead i've set up [cert-manager](https://github.com/jetstack/cert-manager)
which can store certs in a k8s secret and mapped those into traefik

## grpc-web

so far no explicit support for grpc-web,
your server still needs to understand grpc-web on its own
at least the reverse proxy doesn't break anything
