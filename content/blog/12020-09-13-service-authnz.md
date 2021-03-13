---
description: authn/z by the service or network
title: service authn/z
---

### scenarios

#### _service_ authn authz

The service handles authn and authz by itself.
The network just routes everything to it.

Advantages:

- can be deployed standalone
- dumb network

Disadvantages:

- disparate identities across services
- service needs state about authn / authz

#### _service_ authn authz _external_ identity

The service handles authn and authz,
but delegates authn to an external identity provider.
The network just routes everything to it.

Advantages:

- can be deployed standalone
- dumb network
- shared identity across services

Disadvantages:

- service needs state about authn / authz

#### _network_ authn _service_ authz

The network handles authn,
blocking/redirecting requests to be authenticated first.
The service trusts the identity provided by the network,
ex: as headers/jwt.
The service decides on authz based on provided identity.

Advantages:

- shared identity across services
- service doesn't need state about identity management

Disadvantages:

- always need to run beind proxy/mesh
- extra service to handle authn
- service still needs to map identity -> authz

#### _network_ authn authz

The network handles authn and also authz,
deciding if requests should pass to service.
The service trusts all requests as valid.

Advantages:

- shared identity across services
- shared authz across services
- service doesn't need any state about authn authz
- unauthorized requests can be stopped at the edge

Disadvantages:

- always need to run behind proxy/mesh
- extra service(s) to handle authn authz
- pushes complexity into network
- moderately complex service will likely still need some knowledge of authz
  (less of a problem if it can be stateless)
- network is likely only able to do coarse grained authz
