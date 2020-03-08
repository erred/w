---
description: benchmarking building containers in cloud build
title: benchmark cloudbuild
---
### Benchmark building containers

for my use case: small Go microservice repos with multistage build

_conclusion_: use `docker build`

#### test case

reference repo: [authed](https://github.com/seankhliao/authed),
**100** lines of code, +vendored `grpc, grpcweb, firebase`

cache == nothing changed _.dockerignore_

#### kaniko no cache multistage

89s, 85s

#### kaniko cache multistage

93s, 85s

#### docker build no cache multistage

60s, 57s

#### docker build cache multistage

58s, 56s

#### buildah no cache multistage

140s, 138s

#### img no cache multistage

130s, 120s