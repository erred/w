---
title: cicd and gitops
description: both trendy, but do they mix?
---

### _ci/cd_

Conitnuos Integration:
roughly translates to build early, build often,
and never deviate far from mainline,
on the idea the many small changes are easier to swallow
than a single big one.

Continuous Deployment:
You have a steady stream of build artifacts, now what?
You deploy and run them often,
hopefully you spot breaking changes early and fix them sooner rather than later.

GitOps:
There's _desired state_ and _reality_,
and gitops is about tracking _desired state_ in git
so you have a clear history and it can be reviewed before being rolled out.
Usually/hopefully you have you have automated tooling to do the rollout part,
so humans don't get the chance to make mistakes.

#### _implementation_ questions

Theory sounds all nice and good,
but how do you implement this?

##### _linear_ single stage

This is easy to understand:
as part of your CI pipeline,
you run the tasks:

- build
- push artifacts
- deploy

to redeploy, an existing release,
either rerun part or all of the pipeline (dummy commits anyone?)
or do it manually.
A critical weak point when it comes to gitops with this setup is
there's no good way for automation to record the deploy version into git for a single repo/branch.
If you use commit hash (or any other generated identifier),
modifying the deployment manifests to record it changes the identfier.

Advantages: easy to understand and implement.

Disadvantages: sometime after deploy _reality_ may go out of sync with what you pushed,
less direct control over deploy timing,
gitops versioning problem.

##### _linear_ 2 stage

Like the previous version but split into 2 stages:

- build
- push artifacts
- update manifests

and

- deploy manifests

This is usually done with 2 repos,
it's your choice of placing the canonical version of deployment manifests with code
and generating the entire thing into a separate repo
or having the canonical version in the separate repo and just updating the artifact id every build.

The deploy stage is just another pipeline.

Advantages: this solves most of the problems from linear single stage:
no more versioning problems
and you have better control over deploy timing.

Disadvantages: _desired state_ and _reality_ can still go out of sync.

##### _linear/declarative_ mix

This builds on the 2 stage setup,
but instead of having a pipeline for the second stage,
you have a controller watch both desired state and reality
and take the actions (redeploying) to keep them in sync.
No more configuration drift!

I've yet to see a declarative approach for the build stage,
but i guess it's not really necessary as artifacts don't drift?
