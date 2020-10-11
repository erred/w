---
descriptions: policy types and enforcement
title: k8s policy enforcement
---

### _Policy_ Enforcemnt

You have _ideas_ about what connections/requests should go through.
Those are policies.
These are some things that can define/enforce those policies.

#### _NetworkPolicy_

[NetworkPolicy][networkpolicy] is the standard API defined by kubernetes
for networking plugins (CNI) to enforce.
It is the lowest common denominator L3/L4,
supporting IP range based allow/deny
and TCP/UDP port based allow/deny.

#### _Cilium_

[Cilium][cilium] is a CNI that can enforce NetworkPolicy,
with its main claim to fame being it is eBPF based.
Additionally extends support to L7 through injecting an envoy sidecar
with policies defined in [CiliumNetworkPolicy][ciliumnetworkpolicy].
Also has support for L7 policies (HTTP, Kafka, DNS) through injecting envoy.
HTTP supports Host / Path / Method / Headers (header exact match).
Identity aware in the context of Cilium refers to service identity
through k8s names.

#### _OpenPolicyAgent_

[OpenPolicyAgent][opa] is a generic policy rules engine
that can be embedded or run as a standalone service.
Other components (ex API gateways, k8s API) can be configured to talk to it
to determine authz.

#### _Istio_

[Istio][istio] is a service mesh built on envoy.
It uses [AuthorizationPolicy][istioauthorization] to manage TCP / HTTP traffic policies
for both North-South and East-West traffic.

[networkpolicy]: https://kubernetes.io/docs/concepts/services-networking/network-policies/
[cilium]: https://cilium.io/
[ciliumnetworkpolicy]: https://docs.cilium.io/en/v1.8/policy/
[istio]: https://istio.io/
[istioauthorization]: https://istio.io/latest/docs/tasks/security/authorization/
