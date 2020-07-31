---
description: notes on reverse proxies
title: reverse proxies
---

### _reverse_ proxies

_tldr_ use envoy

For some reason or another, you need a reverse proxy.
you have options:

| proxy              | language | config                       | TCP | UDP | http2        | http3        | grpc |
| ------------------ | -------- | ---------------------------- | --- | --- | ------------ | ------------ | ---- |
| _[Envoy][envoy]_   | C++      | yaml / xDS dynamic           | yes | yes | front / back | front / back | yes  |
| [HAProxy][haproxy] | C        | custom text / custom dynamic | yes | no  | front / back | planned?     | yes  |
| [Nginx][nginx]     | C        | custom text                  | yes | yes | front        | front        | yes? |
| [Traefik][traefik] | Go       | yaml / toml / multi dynamic  | yes | yes | front / back | planned?     | yes  |

#### _Notes_

- envoy can be configured statically for basic setups
- envoy is used as a base for more complex dynamic setups, eg k8s ingress
- haproxy is ehhhh on feature support
- nginx is popular but feels compromised feature wise by enterprise edition
- traefik is simple to setup for standalone operations
- traefik include native support for many dynamic config situations

[envoy]: https://www.envoyproxy.io/
[haproxy]: http://www.haproxy.org/
[nginx]: https://www.nginx.com/
[traefik]: https://containo.us/traefik/

#### _Other_

- [Apache HTTPD][apache] is a pure http (1/2) server, not really well suited for proxying
- [MetalLB][metallb] for bare metal deployments + k8s, won't work in clouds

[apache]: https://httpd.apache.org/
[metallb]: https://metallb.universe.tf/
