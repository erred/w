---
description: the foundations of a modern "cloud native" stack
title: foundations of cloud native
---

### _cloud_ native

Technologies and hype cycle come and go.
The bleeding edge today might be the forgotten rusty blade tomorrow.
So it's important to choose the right
(ie. going to be well supported) tools for building infra,
which will likely need to run for quite a while.

Nobody has a crystal ball
(if you do please share),
but the community (people/companies with marketing power)
is moving in certain directions.

see [CNCF landscape](https://landscape.cncf.io/)

#### _core_ technologies

These things are unlikely to get replaced anytime soon

- kubernetes:
  operating system for the cloud,
  and the shared common denominator of infrastructure as code.
- prometheus:
  pull based metrics collection and storage
- grafana:
  graphing metrics

#### _emerging_

Things that already form core components but still need a bit more convincing
that they are essential

- envoy: L4/L7 API driven proxy
  - space needs more consolidation
- cert-manager: TLS certificate thing
  - needs more stability
- istio: service mesh
  - needs more ergonomics
- open policy agent: policy evaluation engine
  - needs more demonstrated use
- coredns
  - needs more demonstrated use outside k8s
- etcd
  - needs more demonstrated use outside k8s
