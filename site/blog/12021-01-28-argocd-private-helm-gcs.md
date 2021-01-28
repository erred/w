---
description: does it work?
title: argocd private helm gcs
---

### _argocd_

We've settled on [Argo CD](https://argoproj.github.io/argo-cd/)
for our continuous deployments.
Some other parts of our stack:
[helm](https://helm.sh/) for packaging kubernetes manifests
and [helm-gcs](https://github.com/viglesiasce/helm-gcs)
(bash version) for a private helm repository.

#### _doesn't_ work

So a colleague was tasked with setting it up.
Everything worked more or less as expected,
well almost everything,
they said the helm-gcs plugin didn't work.
Sprint planning comes around and now there's a ticket for solving the issue
with 3 options: fixing it as an upstream bug (linking to [#4439](https://github.com/argoproj/argo-cd/issues/4439)),
using an alternate chart repository,
and something else i forgot.
Side note: an option would also be to expose the bucket as a static site.
Noises were made about it being a difficult problem,
I ask if they tried the other [helm-gcs](https://github.com/hayorov/helm-gcs)
and it was promptly dumped on me.

#### _debugging_

Ah, what do we do? Try to reproduce the problem of course.

- create new image with helm-gcs and redeploy
  - `gsutil` not installed :facepalm:
- crib the [gcloud sdk](https://cloud.google.com/sdk/docs/install) install steps from another dockerfile and try again
  - issue reproduced

What next? Try to see what's happening.

- tail the logs
  - .... failed with `found in Chart.yaml, but missing in charts/ directory`
  - hmm, dependencies not being installed? weird
- add `helm dependency build` to our helm wrapper script (we need it for [helm-secrets](https://github.com/jkroepke/helm-secrets))
  - it works
  - but only once, then it fails with `Unable to move current charts to tmp dir: rename charts tmpcharts: file exists`
  - google the error, `tmpcharts` is from when helm is killed
- tail the logs some more
  - notice that argocd is actually calling `helm dependency build` after the first failure
  - apparently it's an optimization so the first error is expected
  - remove hack in wrapper script
  - notice the real error from is timing out after 1m30s

That was a dead end, time to get my hands dirty

- exec into the container and find the checked out dir
  - run `time helm dependency build`
  - `1m33s` so borderline so it sometimes works and sometimes (mostly) doesn't
  - why is it taking so long though, don't remember it being so slow
  - run it again, think of something else, `Ctrl-C`
  - python stacktrace dumped on my screen, oh i get it

The helm-gcs we use is a bash script that wraps `gsutil` which is itself written in python.
Argo itself is written in Go, so the resources were sized accordingly,
request/limit of `25/50m cpu` and `64/128Mi memory`.
Python is hungry and fat.

- up the resources to `100/200m cpu` and `128/256Mi memory`
  - try again clearing out all hacks
  - it works, consistently, also much faster

I believe my original proposition of using the other helm-gcs (written in Go)
might also have worked, but wouldn't have uncovered this issue of resource starvation.
