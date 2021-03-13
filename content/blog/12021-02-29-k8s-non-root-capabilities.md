---
description: using capabilities in kubernetes when you're not root
title: k8s non root capabilities
---

### _non_ root

By default, containers are root in their own little sandbox,
which is nice from a don't-need-to-think-about-permissions
perspective, but a bit less nice if you're concerned about privilege escalation.

So you have the ability to run as not-root,
either by setting `USER xxx` in a `Dockerfile`
or at runtime.
But sometimes you need to do special things, like binding to port `80`.
Linux gives us [capabilities](https://man.archlinux.org/man/capabilities.7.en),
granting little slices of elevated permissions.

#### _docker_ / kubernetes

[Unfortunately](https://github.com/kubernetes/kubernetes/issues/56374)
for us, it doesn't work just by granting the capability to everything in the container.
So you actually need to use [setcap](https://man.archlinux.org/man/setcap.8.en)
to set the permissions you want on the binary
(sometimes necessitating a custom image build) before the capabilities will work.

Dockerfile for a simple Go http server

```Dockerfile
FROM golang:1.16-alpine AS build
WORKDIR /workspace

# static binary for scratch container
ENV CGO_ENABLED=0

# get the setcap command
RUN apk add --update --no-cache libcap

# main.go:
#     package main
#
#     import (
#             "log"
#             "net/http"
#     )
#
#     func main() {
#             log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
#                     rw.Write([]byte("ok"))
#             })))
#     }
COPY go.mod main.go .
RUN go build -o /bin/app .

# set net_bind both [e]ffective and [p]ermitted
RUN setcap cap_net_bind_service+ep /bin/app


FROM scratch
COPY --from=build /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
```

k8s Pod (can also be combined with PodSecurityPolicy)

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: http
spec:
  # run as not root
  securityContext:
    runAsGroup: 65535
    runAsNonRoot: true
    runAsUser: 65535
  containers:
    - name: http
      image: http:latest
      imagePullPolicy: IfNotPresent
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          # set permitted privileges
          add:
            - NET_BIND_SERVICE
          # default drop all
          drop:
            - all
        privileged: false
        readOnlyRootFilesystem: true
```

#### _bonus_

test with KinD

```sh
docker build -t http .
kind create cluster
kind load docker-image http
kubectl apply -f pod.yaml
```
