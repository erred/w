---
description: reflections on a k8s resource labelling initiative at work
title: k8s labelling reflections
---

### _kubernetes_ labels

Kubernetes resources all have a `metadata.labels` field,
allowing you to add key-value pairs to help with tracking / selecting things.

So we discovered Open Policy Agent
[Gatekeeper](https://github.com/open-policy-agent/gatekeeper),
A way for you to force all resources to conform to policies,
otherwise prevent their creation.
One of these policies is enforcing a required set of labels.
To use this, we had to ensure this didn't cause any unintended disruptions,
so off we go adding a bunch of labels to everything.

#### _only_ add labels you will actually use

There are a lot of things that sound nice,
like `owner`, `environment`, `version` labels
on everything.
But after you've wasted 2 weeks adding labels to everything,
are you actually going to use it?
Or are they there just so you can fulfill a policy (that you also wrote)?

#### _bad_ tooling

If you use Helm, you'll know the inconsistency of upstream charts,
not everything uses the
[recommended set](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/)
much less provide you with options to insert labels on everything.

Even then, it very likely has different ideas
on what the common labels mean versus what's actually useful for you.
eg:

- `environment` could just be `app.kubernetes.io/instance`
- `version` could be `app.kubernetes.io/version`

#### _namespace_ per app

Adding labels to everything is almost certainly a giant waste of time,
better to isolate each app in their own namespace and just label that.
