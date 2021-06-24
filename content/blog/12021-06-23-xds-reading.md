---
title: xDS reading
description: background reading on xDS protocol
---

### _xDS_

xDS is a _set_ of discovery services (hence the name _x_ Discovery Service or _xDS_)
used to convey routing information for a service mesh or load balancing.
I'm aware of 2 main users:
[Envoy](https://www.envoyproxy.io/) (who primarily designed it) and
[gRPC](https://grpc.io/).

#### _background_ info

Envoy's blog on [universal data plane](https://blog.envoyproxy.io/the-universal-data-plane-api-d15cec7a)
is probably the best intro into why/how the data and control plane are separated
and how the api came to be.
The [gRPC xds proposal](https://github.com/grpc/proposal/blob/master/A27-xds-global-load-balancing.md)
introduces a lot of the concepts more cleanly than the envoy docs.

As a set of APIs defined in protobuf/gRPC,
The definitions can be found in the main envoy repo under
[envoy/api/envoy/service](https://github.com/envoyproxy/envoy/tree/main/api/envoy/service)
mirrored to the [envoy/data-plane-api](https://github.com/envoyproxy/data-plane-api) repo
for consumption.
[cncf/xds](https://github.com/cncf/xds) appears to be a future home for it?
There are also reference control plane servers, ex: [go-control-plane](https://github.com/envoyproxy/go-control-plane).


#### _glossary_

- Listener: Local L4/L7 listener + filters for proxy
- Route: HTTP routing table
- Cluster: Upstream to forward requests to
- Endpoint: Backing hosts for upstreams

- ADS: Aggregated Discovery Service
- CDS: Cluster Discovery Service
- CSDS: Client Status Discovery Service
- EDS: Endpoint Discovery Service
- LDS: Listener Discovery Service
- RDS: Route Discoery Service
- RTDS: Runtime Discovery Service
- SDS: Secret Discovery Service
- SotW: State of the World
- SRDS: Scoped Route Discovery Service
- VHDS: Virtual Host Discovery Service
- UDPA: Universal Data Plan API
- DPLB: Data Plane Load Balancer
- CPLB: Control Plane Load Balancer
- LRS: Load Reporting Service
- ORCA: Open Request Cost Aggregation
