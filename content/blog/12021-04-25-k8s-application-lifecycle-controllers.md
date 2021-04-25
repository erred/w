---
title: k8s application lifecycle controllers
description: something fancier than a plain Deployment
---

### _application_ lifecycle

There's so much more to an application than just its straight runtime environment (Deployment).
Things like the upgrade/rollout policy, templating multiple similar applications etc...

#### _defining_ application

##### _oam_

[Open Application Model](https://oam.dev/) is a spec for component reuse.
Define `ComponentDefinitions` that takes in parameters
and outputs templated (`CUE`, `Helm`, raw) resources,
and compose / instantiate them with `Applications`.

###### _kubevela_

[KubeVela](https://kubevela.io/) is an implementation of oam.
Update the `Application`
and KubeVela will create/update the underlying components for you.

Also does traffic shifting with `AppDeployment`(?)

##### _kpt_

[kpt](https://googlecontainertools.github.io/kpt/) is like a fancier version of kustomize,
take in any raw upstream source, add setters/substitutors, and apply with a safer(?) command.

#### _deploying_ new versions

##### _knative_ serving

[Knative Serving](https://knative.dev/docs/serving/) has a `Service` type
which wraps routing and pod definition into a single resource,
handling scaling (scale to zero) and a bit of traffic routing.

##### _flux_

[flux](https://toolkit.fluxcd.io/core-concepts/)
watches Git repos / storage buckets / image registries,
runs Helm / Kustomize templating and applies resulting resources.

##### _argo_ cd

[Argo CD](https://argoproj.github.io/argo-cd/)
watches Git repos, runs templating tools
and applies resulting resources.

Like flux, but at a bigger scale
(multicluster, users, workspaces, ui).

##### _kudo_

[KUDO](https://kudo.dev/docs/architecture.html#architecture-diagram)
define your app wth Go templating for resources
and an `operator.yaml` describing steps to take during create / update / delete.
Pass in values through a `parameters.yaml` and apply with a `kubectl kudo` cli.

Basically a fancier version of Helm 2...

#### _rolling_ out changes

These manage traffic routing to new versions of your deployments.

##### _argo_ rollouts

A [Rollout](https://argoproj.github.io/argo-rollouts/features/specification/)
is CRD that is basically a `Deployment` with extra fields to reference `Services`
to manage the upgrade lifecycle.
Editing a live, existing `Rollout`
will result in a new `ReplicaSets` being created,
traffic being routed to it, before final promotion.

##### _flagger_

[Flagger](https://docs.flagger.app/usage/how-it-works) takes a different approach,
instead watching a `Deployment` or `DaemonSet` and creating a "stable" primary clone.
Traffic is routed to the primary clone and temporarily diverted during upgrades.

#### _everything_

why run all these separate components when it could be just one big monolith?

##### _keptn_

[keptn](https://keptn.sh/) takes in a helm chart,
and on trigger (cli, CloudEvents?)
runs through a gauntlet of rollouts / tests before deploying with custom steps.
After deploying, continously monitor and trigger more custom steps to handle problems.
