---
description: reading the service mesh interface spec
title: reading smi spec
---

### _service mesh_ interface

Service meshes are hype.
Maybe, not really, yet.
They're still a lot of overhead
(so many sidecars) but they do provide some valuable features.
Wouldn't it be sad if every mesh implementation defined its own config,
locking you in?

[Service Mesh Interface](https://smi-spec.io/) or SMI
is the proposed spec for standardizing access to service mesh features.
Unfortunately, it is still very much alpha,
and the [spec](https://github.com/servicemeshinterface/smi-spec/blob/main/SPEC_LATEST_STABLE.md)
is little more than the api CRD descriptions.
The metrics API does seem the most fleshed out,
but it does seem it will be a while fore it's fully ready.
