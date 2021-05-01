---
title: normalizing k8s manifests
description: how do you properly compare manifest yamls
---

### _manifests_

Kubernetes manifests most commonly come as YAML,
which can be a problem if you want to compare 2 manifests.
The field order doesn't matter,
there are a bazilion different syntaxes that mean the same,
and quoting issues abound.

So the best way I cam up with was to run everything through
[kustomize](https://kustomize.io/),
it generates consistent, formatted output that can easily be passed to diff.
