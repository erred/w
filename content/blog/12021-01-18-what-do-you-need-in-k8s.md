---
description: k8s, the build-a-bear diy PaaS thing
title: what do you need in k8s
---

### _kubernetes_

You're thinking of running [kubernetes](https://kubernetes.io/).
Let's think this through first.

#### _before_ you start

Where do you want to run it?
In the cloud with a managed control plane:
[GKE](https://cloud.google.com/kubernetes-engine),
[EKS](https://aws.amazon.com/eks/)?
In the cloud with bare VMs?
Or on prem?

If it's managed,
do you use [terraform](https://www.terraform.io/),
use the cloud provider specific CLI / SDK / config management,
or do you have a snowflake cluster where you can use [gardener](https://gardener.cloud/).

If not, you need to at least choose:
a [CNI](https://github.com/containernetworking/cni) provider for networking,
and a [CSI](https://github.com/container-storage-interface/spec) provider for storage.
Which one is probably going to be affected by the features you need later.

Also, you need to decide how you want to spin up / manage your nodes / kubelets.
[kubeadm](https://github.com/kubernetes/kubeadm),
[k3s](https://k3s.io/), ...
or something more specialized, maybe [k0s](https://k0sproject.io/).

#### _config_ management

You have a cluster, now you want to install stuff into it.
Best to decide now which tool you want to use and _standardize_ on it.
[helm](https://helm.sh/) is _the_ package manager, but it's a shitty one,
you're almost certainly better off without it unless you run stock everything.
You can always use raw manifests
but [kustomize](https://kustomize.io/) is a worthwhile layer on top,
it even comes built in with [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/).

Or you just need to be different,
[grafana tanka](https://grafana.com/oss/tanka/) uses jsonnet,
terraform has a [kubernetes](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs) provider.

Whatever you choose, you also have the problem of secrets.
[sealed-secrets](https://github.com/bitnami-labs/sealed-secrets) are a write only solution,
a lot of other ones use [sops](https://github.com/mozilla/sops) under the hood,
like [ksops](https://github.com/viaduct-ai/kustomize-sops).
Or be like me and hack up some dingy solution with [gitattributes filter](https://git-scm.com/docs/gitattributes#_filter)
and [age](https://github.com/FiloSottile/age).

#### _cluster_ basics

It's alive, now what?
If you serve web traffic,
you'll want an [ingress controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)
to give you more than L4 routing,
[ingress-nginx](https://kubernetes.github.io/ingress-nginx/) is default
(not to be confused with [nginx ingress](https://www.nginx.com/products/nginx-ingress-controller/)),
[traefik](https://traefik.io/) is pretty popular if you're not doing weird stuff,
or maybe you already know you need a service mesh,
[istio](https://istio.io/) has the most mindshare,
[linkerd](https://linkerd.io/) is another big one,
others seem more half-hearted.

It's 2021, you need to serve TLS,
[cert-manager](https://cert-manager.io/) is almost mandatory,
even if it is still a pain to manage/upgrade.
If you do the sane thing and give each service its own subdomain,
you'll need to manage your DNS entries too,
[external-dns](https://github.com/kubernetes-sigs/external-dns) can do that for you.

Oh and if you're not in the cloud, you'll need something like [MetalLB](https://metallb.universe.tf/),
hostports, or some other way of exposing your ingress controllers to the outside world,
[envoy](https://www.envoyproxy.io/) with some static config
will probably also work, you're only exposing a single service anyway.

#### _observability_

What next? you can deploy your application now and it'll serve traffic just fine.
But you want to know what it's doing, how well it's performing
and have something to look at when things go wrong.
Say hello to the 3 pillars of oberservability: _metrics_, _logs_, _tracing_.

[prometheus](https://prometheus.io/) is what everyone uses for metrics
(unless you use some hosted thing), it scrapes metrics from various services and stores it.
If you run at some mind boggling scale, you can consider
[thanos](https://thanos.io/) or [cortex](https://cortexmetrics.io/)
for scaling prometheus, but for most people a beefy prometheus (HA pair if you care) is enough.

Along with prometheus, you'll want:
[kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) to expose k8s things to prometheus,
[node exporter](https://github.com/prometheus/node_exporter) for metrics about your host,
[pushgateway](https://github.com/prometheus/pushgateway) if you have things that can't be scraped,
and [alertmanager](https://github.com/prometheus/alertmanager) for alerting.

For logs you'll want [promtail](https://grafana.com/docs/loki/latest/clients/promtail/) to scrape the logs,
and [loki](https://grafana.com/docs/loki/latest/) to store them.

Don't forget the most important part, pretty dashboards!
Do you even have a choice other than [grafana](https://grafana.com/)?

For tracing, well...
until [opentelemetry](https://opentelemetry.io/) manages to stop making breaking api changes every few weeks,
one of [jaeger](https://www.jaegertracing.io/) or [zipkin](https://zipkin.io/) will have to do
(both have committed to eventually converge on opentelemetry).
Maybe you can store in [grafana tempo](https://grafana.com/oss/tempo/)?

If you're fancy, you can run the [opentelemetry collector](https://github.com/open-telemetry/opentelemetry-collector)
to ingest traces and metrics so you can preprocess / filter / reexport them,
you still need to store the data somewhere though, so it doesn't really save you from running the above components.

#### _cicd_

Now... you want a fluid way of getting your application code
(or any of the previous components) into a running cluster.
`git push` should be all you need!

For _CI_ you may want some solution integrated with you code host,
or you want to run something yourself,
[tekton](https://tekton.dev/) is a decent choice.

What about _CD_? What about it? If you believe in GitOps
(code in git describes desired state, controllers make it happen),
[fluxcd](https://toolkit.fluxcd.io/) is available but
[argocd](https://argoproj.github.io/argo-cd/) is more mature.

#### _security_

Ah, better have a good threat model in mind before you start thinking about this.

For you cluster, you'll want to lock down _RBAC_,
and since it's allow only (no deny) [hierarchical namespaces](https://github.com/kubernetes-sigs/multi-tenancy/tree/master/incubator/hnc)
make managing permissions less daunting.

[LimitRange](https://kubernetes.io/docs/concepts/policy/limit-range/) can ensure nothing will exhaust all your resources,
[PodSecurityPolicies](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) limits the permissions of your pods,
and [NetworkPolicies](https://kubernetes.io/docs/concepts/services-networking/network-policies/) limits the cross talk between your services,
though if you use istio you can control it at L7 with [AuthorizationPolicies](https://istio.io/latest/docs/reference/config/security/authorization-policy/).

I'm sure you can dream up of various things you want to enforce on your k8s objects
[OpenPolicyAgent Gatekeeper](https://github.com/open-policy-agent/gatekeeper) makes that possible.

What about protecting your services?
[pomerium](https://pomerium.io/) is a decent solution to implement SSO for multiple services with an external provider, integrating with your ingress controller.
I'm sure you could cook up something similar with the [ory](https://www.ory.sh/) projects.

#### _are_ we done yet

Maybe? Are we ever done?
At least your resume is now filled with exciting new technologies...
