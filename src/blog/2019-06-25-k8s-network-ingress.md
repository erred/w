---
description: reaching your cluster
title: k8s network ingress
---
so you have a kubernetes cluster and now you want to reach it from the outside,
you have several options

### what works

spend money on a hosted LoadBalancer,
use DNS as a ghetto load balancer,
or run your own (edge router outside the cluster on stable ip + some way of service discovery)

#### Service: LoadBalancer

this provisions a cloud hosting provider managed _load balancer_,
it costs \$\$$,
like $18 per month per route on GCP (for the first 5),
you could run a cluster for less
**especially if you run them on preemptible instances, where you need this even more\_**

#### Deplyment: hostPort

this binds to a specified port on the host node,
not too useful if you use preemptible instances and your ips change often

best used if you have a small _stable_ node to run the ingress
and everything else runs on preemptible instances,
still has problems updating:
k8s wants a new version up before if kills the old
but that's impossible because the port is in use

#### tunnelling ingress

like _cloudflare argo_?
basically a load balancer but set up differently

### What doesn't work

#### Service: NodePort

only useful if you don't mind a high port number

#### Service: ExternalIPs

this just tells kube-proxy to allow packets destined to the ips to be routed,
doesn't actually listen on the ports

#### external-dns

[kubernetes-incubator / external-dns](https://github.com/kubernetes-incubator/external-dns)
only points an external DNS name to an external IP created by a LoadBalancer or something else

so no pointing it directly at the node