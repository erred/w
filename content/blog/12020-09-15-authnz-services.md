---
description: services for authnz
title: authnz services
---

### _authnz_

For user requests only.

#### _authn_

Situation:
you have an API gateway (ex traefik)
or service mesh (ex istio)
that can intercept and redirect requests.
You have some existing authz system you want to use,
maybe the same as for your East-West traffic.
You want to redirect requests without a valid session / jwt
to login to a third party identity provider.

Options:

- [Pomerium][pomerium]
- [ORY Kratos][kratos]
- [loginsrv][loginsrv]
- [Auth0][auth0]
- [Keycloak][keycloak]

#### _authn authz_

Situation:
you have an API gateway
that can intercept requests.
You want to redirect requests without a valid session / jwt
to login to a third party identity provider
and enforce some policy on that request at the same time.

Options:

- [Pomerium][pomerium] (internally uses OpenPolicyAgent)
- [ORY Kratos][kratos] + [ORY Keto][keto] (Keto is built on OpenPolicyAgent)
- [Pomerium][pomerium] / [ORY Kratos][kratos] + [OpenPolicyAgent][opa]

#### _identity_ provider

loginserver4, ORY Hydra, ...

#### _gateway_

The reverse proxy / identity aware proxy / thing that intercepts your requests.

traefik, envoy, ory oathkeeper, pomerium, ...

[pomerium]: https://pomerium.io/
[loginsrv]: https://github.com/tarent/loginsrv
[identityserver]: https://github.com/IdentityServer/IdentityServer4
[keycloak]: https://www.keycloak.org/
[opa]: https://www.openpolicyagent.org/
[auth0]: https://auth0.com/
[ory]: https://www.ory.sh/
[oathkeeper]: https://www.ory.sh/oathkeeper
[keto]: https://www.ory.sh/keto
[kratos]: https://www.ory.sh/kratos
