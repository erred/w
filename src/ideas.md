---
description: backlog of ideas
header: <h2><em>ideas</em> backlog</h2> <p>append only list</p>
style: |-
  main h3 {
  margin-top: 25vh;
  }
title: ideas
---

### _Dev_ implementation required

#### _ci:_ tekton webhook interceptor

cel is nice but a bit too limiting,
implement custom interceptor with plugin style

#### _ci:_ get latest version

retrieve the latest version of thrid party X
from a variety of sources (github, arch, other).
update / tag repo and push for ci

#### _go:_ module proxy over IPFS

modify `go` command to retrieve modules over IPFS
maybe works well with IPFS deduplication?

#### _web:_ webauthn

try it

#### _site:_ serve over AMP

this site doesn't need fancy things,
generate and serve over AMP with subdomain & use with cloudflare

#### _site:_ serve with WebPackaging / SXG

serve signed / shareable bundle for site / page

- needs proper SXG cert: expensive

#### _site:_ serve over IPFS

use DNSLink and serve this site over IPFS

#### _site:_ serve over Tor alt-svc

apparently it's not difficult?

#### _archlinux:_ mirror/repo over IPFS

Arch should be flexible enough to handle getting packages over ipfs
with a special `XferCommand`.

#### _archlinux:_ reflector

reimplement reflector in go because NIH

### _Ops_ time and patience needed

#### _ci:_ build on push

remove limitation on tagged commits

#### _ci:_ module proxy for go

run a private module proxy in CI cluster as a local cache

#### _ci:_ container registry

run a private container registry in CI cluster as a cache for kaniko

#### _cd:_ deployment for k8s services

as title

#### _archlinux:_ full disk encryption

add full disk encryption to XPS

### _Write_ wordsmithing wanted

#### _like:_ more go modules

- [auth](https://github.com/avelino/awesome-go#authentication-and-oauth)
- [cli](https://github.com/avelino/awesome-go#standard-cli)
- [config](https://github.com/avelino/awesome-go#configuration)
- [data](https://github.com/avelino/awesome-go#database)
- [tui](https://github.com/avelino/awesome-go#advanced-console-uis)
- [uuid](https://github.com/avelino/awesome-go#uuid)

#### _blog:_ tekton pipelines

write experience report

#### _blog:_ SNE courses review / learned

experiences / things learned

was it worth it?
