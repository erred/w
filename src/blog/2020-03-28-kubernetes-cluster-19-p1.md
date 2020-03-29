---
description: 19th attempt to get a k8s setup working
title: kubernetes cluster 19 part 1
---

### Goals

- kubernetes cluster
- config in git
- tls certificates

#### _tldr_

config: [seankhliao/kluster @ v0.19.0][kluster]

use with:

- `make create-cluster`
- `make decrypt` // or create appropriate secret files
- `kubectl apply -k .` // maybe repeat a few times if things fail

##### cluster / config

- [GKE][gke]
- kustomize

##### certificates

- [cert-manager][certm] and [lets encrypt][lets]

#### cluster

[GKE][gke], but _cheaply_, ok? _1_ zonal [E2][e2] node

- disable http load balancing
- add DNS + wildcard CNAME to point to node

#### config

use `kubectl`'s built in [kustomize][kustomize] support

- [helm 3][helm] might be better than 2 without Tiller, but I want control
- helm is still useful to get a decent starting config `helm install --dry-run --debug name repo/chart > bundle.yaml`
- the only command you'll ever need is `kubectl apply -k .`
- careful with `--prune --all` especially with GKE managed resources(?)
- `configMapGenerator` and `secretGenerator` means config and secrets get their own files (and file types!)

#### tls

[cert-manager][certm] is a royal pain to get running. DO NOT attempt to change anything from their default config

Lets Encrypt Cluster Issuer with Cloudflare DNS Challenge

- cloudflare API token needs `Zone:Zone:Read`, `Zone:DNS:Edit` on all zones(?)

```yaml
apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: le-issuer
spec:
  acme:
    email: admin@example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: le-issuer-account
    solvers:
      - dns01:
          cloudflare:
            email: admin@example.com
            apiTokenSecretRef:
              name: cloudflare
              key: token
```

certificate resource:

```yaml
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: api.seankhliao.com
spec:
  secretName: api-seankhliao-com-tls
  duration: 2160h
  renewBefore: 360h
  dnsNames:
    - api.seankhliao.com
    - "*.api.seankhliao.com"
  issuerRef:
    name: le-issuer
    kind: ClusterIssuer
```

[lets]: https://letsencrypt.org/
[helm]: https://helm.sh/
[kustomize]: https://github.com/kubernetes-sigs/kustomize
[certm]: https://github.com/jetstack/cert-manager
[kluster]: https://github.com/seankhliao/kluster/tree/v0.19.1
[gke]: https://cloud.google.com/kubernetes-engine
[e2]: https://cloud.google.com/blog/products/compute/google-compute-engine-gets-new-e2-vm-machine-types
