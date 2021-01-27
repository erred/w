---
description: what i've done in the last 2 months
title: 4 months at job
---

### _4_ months in

So, following up on [2 months in](/blog/12020-12-03-junior-sre-at-non-tech/),
what more have I done?

#### _people_ and process

I have a better grasp of what we manage (a lot of stuff)
and a tiny bit of office politics (sigh).
I can start recognizing people,
but only as names though MS Teams or Github or email
(all microsoft products...).

Within our team, I also have a better idea of whose input to lend more weight to,
and whose to question.
Also started doing code reviews,
if only because I feel sorry
for regularly sending out PRs that touch 50-100 files at once.
Yes I try to be detailed without being nitpicky.

#### _work_ and tech

Finally finished up the boring work of upgrades (lots of third party apps)
and migrations (ansible, helm2 -> helm3).
Added some more migrations like an elastic cluster
(that I don't think I did properly), because you can't have too many, right?
Oh, and some code cleanup and standardization while everyone was away
(formatter, unified build/deploy commands with `make`).

After a 2 month fruitless discussion on CI/CD
I can finally push my preferred direction solution (split CI/CD + GitOps).
Started dumping ideas on the team for automation,
among the ones we investigated and are likely to adopt:

- update dependencies with [renovate](https://github.com/renovatebot/renovate)
- continuous integration builds with [github actions](https://github.com/features/actions)
- gitops style continuous deployments with [argocd](https://argoproj.github.io/argo-cd/)

Security is also an ongoing focus:
locking down permissions, having a private copy of all dependencies,
rotating credentials, futzing around with iam,
enabling security features in k8s (
[shielded nodes](https://cloud.google.com/kubernetes-engine/docs/how-to/shielded-gke-nodes)
[cos_containerd](https://cloud.google.com/kubernetes-engine/docs/concepts/using-containerd)
) and possibly in the future (
[hierarchical namespaces](https://github.com/kubernetes-sigs/multi-tenancy/tree/master/incubator/hnc)
[workload identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity)
[network policy](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
[cilium](https://cilium.io/)
) as well as my current ongoing time sink:
[cloudflare bot management](https://www.cloudflare.com/products/bot-management/)

Observability work is coming... soon... hopefully.
I have ideas, but I do want to clear out this set of changes first.
