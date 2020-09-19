---
description: dev loop for go projects in k8s
title: k8s go dev loop
---

### _dev_ loop

#### _manual_

Need a local dev env + registry + a cluster

1. write dockerfile
2. write k8s manifest
3. write go code
4. `docker build && docker push`
5. `kubectl apply / set`
6. go to 3

#### _ci_

Need partial local env + CI system + registry + cluster

1. write dockerfile
2. write k8s manifest
3. write CI config to build + push + deploy
4. write go code
5. `git commit && git push`
6. go to 3

#### _ko_

[ko](https://github.com/google/ko)
is specialized for go,
automagically generate a dockerfile,
basically automates the manual loop.

Need a local dev env + registry + a cluster

1. write k8s manifest (replace image with go import path)
2. write go code
3. `ko apply`
4. go to 2 (or run with `--watch`)

#### _skaffold_

[skaffold](https://github.com/googlecontainerTools/skaffold)
is a more general purpose version of ko,
for more languages, supporting more flows,
such as CI.

Supports rendering final manifests for final deploy

Need: whichever flow you're automating

1. write dockerfile
2. write k8s manifest
3. write go code
4. `skaffold run`
5. go to 3 (or run `skaffold dev`)
