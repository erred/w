---
title: managing k8s manifests
description: too much yaml
---

### _k8s_ manifests

Kubernetes manifests are great, defining a common, declarative,
API-first approach to managing resources.
Just too bad that the canonical format is
[YAML](https://yaml.org/) with all its oddities, especially around types.
And the fact that there isn't much to help you deal with repitition
of same or similar blocks.

So what other options do we have?

#### _helm_

Ah, [Helm](https://helm.sh/), CNCF Graduated,
and the closest thing to a standard package manager for kubernetes.
What you get is using [Go text/template](https://pkg.go.dev/text/template)
to generate YAML (not even structured data),
feeding in variables through another YAML file,
with all the typing problems that comes with it.
It sort-of helps with duplication in a single chart,
or if you have the stomach for it,
a single chart that you can cram all your microservices into.
It's a (not very smart) templating engine with a deployment tool on top,
and your customization points are limited to what the chart authors came up with,
and if you use upstream charts, they're all of varying quality.
you have to deal with varing levels of quality from upstream charts.

#### _kustomize_

A step up from raw kubernetes manifests,
you get some common functionality, like label and annotations for everything,
in a short form with patches and transformers for a fully customizable output.
It does require more sprawl in directory and file setups but it can manage your entire cluster.
You're still writing YAML though... and the correct docs for it are a bit hard to find (hint: last one):

- [kustomize.io](https://kustomize.io/)
- [kubernetes-sigs.github.io/kustomize](https://kubernetes-sigs.github.io/kustomize/)
- [kubectl.docs.kubernetes.io](https://kubectl.docs.kubernetes.io/references/kustomize/)

#### terraform _kubernetes_ provider

[Kubernetes provider](https://registry.terraform.io/providers/hashicorp/kubernetes/latest):
you write the same manifest, but in HCL / terraform.
Good: less YAML to worry about. Bad: HCL is more verbose in lines-of-code (screen real estate),
blocks don't have a good story for reuse in variables, and you're limited to the resource versions they've defined.

#### terraform kubernetes _alpha_ provider

[Kubernetes Alpha provider](https://registry.terraform.io/providers/hashicorp/kubernetes-alpha/latest):
write any k8s resource in HCL / terraform,
which sort of makes you question why, other than you really hate YAML.

#### grafana _tanka_

[Tanka](https://tanka.dev/), write your manifests in [json(net)](https://jsonnet.org), it's not as bad as it sounds,
with some functions and templating blocks and overrides,
but it is a bit exotic?
If you're fully into the Grafana ecosystem,
jsonnet is probably really worth it (also for Grafana dashboards).

#### _cue_

[CUE](https://cuelang.org/), more a language than a fully featured tool,
the core idea seems to be having a schema you can validate / generate against.
The [comparison sections](https://cuelang.org/docs/usecases/configuration/#comparisons) in its usecases
go into more detail,
but in general, probably not as good for managing config sprawl.
