---
title: skaffold kaniko docker caching
description: the convoluted setup i go through to get proper caching
---

### _skaffold_

[skaffold](https://skaffold.dev/) is a great dev tool to detect code changes
and manage the build & deploy phases.
I use it with [kaniko](https://github.com/GoogleContainerTools/kaniko)
to build images in my cluster.
Unfortunately, this comes with a set of problems: caching.



#### _baseline_

Here's our starting setup, using a [random project](https://github.com/seankhliao/feed-agg)
I had lying around.
A standard multistage dockerfile and skaffold pipeline.

```Dockerfile
FROM golang:alpine AS build
ENV CGO_ENABLED=0
WORKDIR /workspace
COPY . .
RUN go build -trimpath -ldflags='-s -w' -o /bin/feed-agg

FROM gcr.io/distroless/static
COPY conf.yaml /etc/feed-agg/conf.yaml
COPY --from=build /bin/feed-agg /bin/feed-agg
ENTRYPOINT ["/bin/feed-agg"]
```

```yaml
apiVersion: skaffold/v2beta16
kind: Config
metadata:
  name: feed-agg
build:
  artifacts:
  - image: europe-north1-docker.pkg.dev/com-seankhliao/kluster/feed-agg
    kaniko:
      reproducible: true
      singleSnapshot: true
      skipUnusedStages: true
      useNewRun: true
      image: gcr.io/kaniko-project/executor:latest
  cluster:
    pullSecretName: kaniko-secret
    pullSecretPath: kaniko-secret
    namespace: skaffold
deploy:
  kubeContext: kind-cluster30
  kustomize:
    paths:
      - kustomize/overlays/cluster30
```

There are a few problems with this:

- it downloads dependencies from a remote proxy every time
- it rebuilds everything from scratch every time

I could split out `COPY go.mod go.sum ./` and `RUN go mod download`
and use docker's layering to cache the dependency download,
but that's meh: it changes every time any dependency changes
and bloats the image repo.

#### _private_ proxy

Instead of downloading from som remote proxy,
we could instead use a local-ish one.
I'm running [athens](https://gomods.io/),
and we can conditionally pass in the proxy using build args,
preserving the defaults when we build elsewhere.

_note:_ `ARG` needs to be specified after `FROM` to take effect inside the image.

This is only very marginally faster than using the public proxy
but does solve one other problem:
private dependencies.
This way no credentials need to be passed to the build image,
since access is implicit.

```Dockerfile
FROM golang:alpine AS build
ARG GOPROXY=https://proxy.golang.org,direct
...
```

```yaml
apiVersion: skaffold/v2beta16
kind: Config
metadata:
  name: feed-agg
build:
  artifacts:
  - image: europe-north1-docker.pkg.dev/com-seankhliao/kluster/feed-agg
    kaniko:
      ...
      buildArgs:
        GOPROXY: http://athens.athens.svc.cluster.local
```

#### _mounted_ caches

To really speed things up,
we need to make use of `go`'s native caching,
namely its local filesystem build and module caches.
We can use `kaniko`'s `--whitelist-var-run`, excluding it from the build image
to mount the caches.

_note:_ I'm using `hostPath` because I'm running a single node
using [kind](https://kind.sigs.k8s.io/)
and couldn't be bothered to set a CSI provider that could provision
ReadWriteMany persistent volumes.

```Dockerfile
FROM golang:alpine AS build
ARG GOMODCACHE=/go/pkg/mod
ARG GOCACHE=/root/.cache/go-build
...
```

```yaml
apiVersion: skaffold/v2beta16
kind: Config
metadata:
  name: feed-agg
build:
  artifacts:
  - image: europe-north1-docker.pkg.dev/com-seankhliao/kluster/feed-agg
    kaniko:
      ...
      buildArgs:
        GOCACHE: /var/run/gobuildcache
        GOMODCACHE: /var/run/gomodcache
      volumeMounts:
        - name: modcache
          mountPath: /var/run/gomodcache
        - name: buildcache
          mountPath: /var/run/gobuildcache
  cluster:
    ...
    volumes:
      - name: modcache
        hostPath:
          path: /opt/kind/cluster30/kaniko-gomodcache
      - name: buildcache
        hostPath:
          path: /opt/kind/cluster30/kaniko-gobuildcache
```

#### _final_ setup

So here's everything,
taking an initial build from 149sec down to 25sec.

```Dockerfile
FROM golang:alpine AS build
ARG CGO_ENABLED=0
ARG GOPROXY=https://proxy.golang.org,direct
ARG GOMODCACHE=/go/pkg/mod
ARG GOCACHE=/root/.cache/go-build
WORKDIR /workspace
COPY . .
RUN go build -trimpath -ldflags='-s -w' -o /bin/feed-agg

FROM gcr.io/distroless/static
COPY conf.yaml /etc/feed-agg/conf.yaml
COPY --from=build /bin/feed-agg /bin/feed-agg
ENTRYPOINT ["/bin/feed-agg"]
```

```yaml
apiVersion: skaffold/v2beta16
kind: Config
metadata:
  name: feed-agg
build:
  artifacts:
  - image: europe-north1-docker.pkg.dev/com-seankhliao/kluster/feed-agg
    kaniko:
      reproducible: true
      singleSnapshot: true
      skipUnusedStages: true
      useNewRun: true
      whitelistVarRun: true
      image: gcr.io/kaniko-project/executor:latest
      registryMirror: mirror.gcr.io
      buildArgs:
        GOPROXY: http://athens.athens.svc.cluster.local
        GOCACHE: /var/run/gobuildcache
        GOMODCACHE: /var/run/gomodcache
      volumeMounts:
        - name: modcache
          mountPath: /var/run/gomodcache
        - name: buildcache
          mountPath: /var/run/gobuildcache
  cluster:
    pullSecretName: kaniko-secret
    pullSecretPath: kaniko-secret
    namespace: skaffold
    volumes:
      - name: modcache
        hostPath:
          path: /opt/kind/cluster30/kaniko-gomodcache
      - name: buildcache
        hostPath:
          path: /opt/kind/cluster30/kaniko-gobuildcache
deploy:
  kubeContext: kind-cluster30
  kustomize:
    paths:
      - kustomize/overlays/cluster30
```
