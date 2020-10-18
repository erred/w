---
description: wading through the growing field of zero trust networking products
title: zero trust networking products
---

### _beyondcorp_

[BeyondCorp](https://www.beyondcorp.com/)
the zero trust security framework developed by Google
to shift the security perimeter to individual people and devices.

Part of that work goes into networking,
creating finer grained tunnels and access controls than what VPNs offer.

#### _overview_

| product                                 | cloud / self  | protocols                    | notes                                                                            |
| --------------------------------------- | ------------- | ---------------------------- | -------------------------------------------------------------------------------- |
| [Tailscale][tailscale]                  | cloud / self  | IP                           | p2p wireguard tunnel setp, technically a VPN                                     |
| [Google IAP][google]                    | cloud hosted  | TCP                          | local proxy, tunnels TCP over HTTPS                                              |
| [Cloudflare Access][cloudflare]         | cloud hosted  | TCP                          | local proxy, tunnels TCP                                                         |
| [Hashicorp Boundary][boundary]          | self hosted   | TCP                          | local proxy, tunnels TCP                                                         |
| [AWS Worklink][aws]                     | cloud hosted  | HTTPS                        | remote desktop/browser?                                                          |
| [Oauth2-proxy][oauth2proxy]             | self hosted   | HTTPS                        | reverse proxy for HTTPS                                                          |
| [yahoo/athenz][athenz]                  | self hosted   | HTTPS                        | reverse proxy for HTTPS                                                          |
| [Pomerium][pomerium]                    | self hosted   | HTTPS                        | reverse proxy for HTTPS / ext auth endpoint                                      |
| [Azure App Proxy][azure]                | cloud hosted  | HTTPS / RDP                  | reverse proxy for HTTPS / RDP?                                                   |
| [Duo Network Gateway][duo]              | self hosted   | HTTPS / SSH                  | reverse proxy for HTTPS / SSH                                                    |
| [Okta Access Gateway][okta]             | self hosted   | HTTPS / SSH                  | reverse proxy for HTTPS / SSH (Adv.Server Access)                                |
| [Trasa][trasa]                          | self hosted   | HTTPS / SSH / RDP            | reverse proxy for HTTPS / SSH / RDP / ext auth endpoint                          |
| [strongDM][strongdm]                    | cloud hosted? | SSH / RDP / k8s? / databases | local proxy, tunnels TCP over TLS, extra support for SSH / RDP / K8s / databases |
| [Gravitational Teleport][gravitational] | self hosted   | SSH / k8s?                   | reverse proxy for HTTPS / some k8s specific support?                             |

other stuff

| product                | cloud / self | protocols | notes                                    |
| ---------------------- | ------------ | --------- | ---------------------------------------- |
| [Smallstep][smallstep] | cloud hosted | SSH       | SSO for SSH (issues SSH certs on demand) |

[aws]: https://aws.amazon.com/worklink/
[azure]: https://docs.microsoft.com/en-us/azure/active-directory/manage-apps/application-proxy
[boundary]: https://www.boundaryproject.io/
[cloudflare]: https://www.cloudflare.com/teams/access/
[duo]: https://duo.com/docs/dng
[google]: https://cloud.google.com/iap
[okta]: https://www.okta.com/products/access-gateway/
[pomerium]: https://pomerium.io/
[gravitational]: https://gravitational.com/teleport/
[sshcom]: https://www.ssh.com/products/privx/
[athenz]: https://github.com/yahoo/athenz
[nassh]: https://github.com/zyclonite/nassh-relay
[oauth2proxy]: https://github.com/oauth2-proxy/oauth2-proxy
[smallstep]: https://smallstep.com/
[strongdm]: https://www.strongdm.com/
[tailscale]: https://tailscale.com/
[trasa]: https://www.trasa.io/
[zscaler]: https://www.zscaler.com/products/zscaler-private-access
