---
description: cutting through more marketing crap
title: more cicd systems
---

### _ci/cd_

The [landscape](https://landscape.cncf.io/category=continuous-integration-delivery&format=card-mode&grouping=category)
shows more projects,
time to read through all their marketing crap (again),
and find new toys to play with.

#### _argo_

You could build a full ci/cd pipeline from all this:
Argo Events to trigger builds in Argo Workflows,
Argo Workflows updating a shared manifest,
Argo CD triggering new deployments,
new deployments go through Argo Rollouts.

Workflows is very barebones when it comes to CI though,
really is generic workflow thing.

##### _workflows_

[argo workflows](https://argoproj.github.io/argo/)

A DAG of steps (k8s pods?) triggered by cli

##### _cd_

[argo cd](https://argoproj.github.io/argo-cd/)

Reconcile git repo state of manifests with cluster state.

##### _rollouts_

[argo rollouts](https://argoproj.github.io/argo-rollouts/)

Replace deployment resource with rollout resource.
Will do canary / blue-green on new versions

##### _events_

[argo events](https://argoproj.github.io/argo-events/)

Takes / listens for events, triggers requests / workflows / rollouts

#### _buildkite_

[buildkite](https://buildkite.com/)

Managed control plane + self hosted runners (an executable you run).
Parallel by default,
allows manual input step,
steps are just scripts/commands you run, no isolation from host.

#### _drone_

[drone](https://www.drone.io/)

Self hosted control plane + runners (ex a docker image).
Each step is a docker container / k8s pod / local exec / ssh session.
Volumes keep state between each step.

#### _flagger_

[flagger](https://flagger.app/)

Controller for adding a canary phase between deployment rollouts.
Add canary resource specifying a deployment,
and new changes to the deployment will go through canary stages first
before being promoted.

#### _flux_

[flux cd](https://docs.fluxcd.io/en/latest/)

Operator for reconciling a declaritive k8s state in a repo with actual state.
Can auto update images based on registry metadata (and sync back repo state).

##### _flux_ toolkit

[toolkit](https://toolkit.fluxcd.io/guides/helmreleases/) aka flux v2

Set of controllers that watch repositories for changes and reconcile
state with helm / kustomize

#### _gocd_

[gocd](https://www.gocd.org/)

Self hosted control plane + runners.
Nothing too interesting(?)

#### _spinnaker_

[spinnaker](http://www.spinnaker.io/)

Looks like a very complex, UI centered way of managing deployments.

#### _keptn_

[keptn](https://keptn.sh/)

Maybe it's something like flagger but manages more of the lifecycle
and provides automated actions such as rollbacks?
Architecture looks overly complex,
docs aren't clear especially on the triggers for deployment,
uses mongodb.

#### _tekton_

[tekton](https://tekton.dev/)

Similar in scope to Argo Workflows + Events,
Pipelines are generic DAG workflows,
Triggers trigger Pipelines.
