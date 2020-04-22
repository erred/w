---
description: backlog of ideas
header: <h2><em>ideas</em> backlog</h2> <p>append only list</p>
style: |-
  main h3 {
  margin-top: 25vh;
  }
title: ideas
---

<!-- markdownlint-disable MD001 -->

### _Dev_ implementation required

Ideas that need code to be written,
may or may not see the light of day

#### _cicd:_ Continuos Integration Pipeline

###### _cicd:_ tekton webhook interceptor

cel is nice but a bit too limiting,
implement custom interceptor with plugin style

###### _cicd:_ get latest version

retrieve the latest version of thrid party X
from a variety of sources (github, arch, other).
update / tag repo and push for ci

partially done for `ci-*` projects on github (with actions)

#### _go:_ go core things

###### _go:_ module proxy over IPFS

modify `go` command to retrieve modules over IPFS
maybe works well with IPFS deduplication?

#### _web:_ webdev things

###### _web:_ webauthn

try it

#### _site:_ this website

###### _site:_ serve over AMP

this site doesn't need fancy things,
generate and serve over AMP with subdomain & use with cloudflare

###### _site:_ serve with WebPackaging / SXG

serve signed / shareable bundle for site / page

- needs proper SXG cert: expensive

###### _site:_ serve over IPFS

use DNSLink and serve this site over IPFS

###### _site:_ serve over Tor alt-svc

apparently it's not difficult?

#### _archlinux:_ Arch Linux Infra

###### _archlinux:_ mirror/repo over IPFS

Arch should be flexible enough to handle getting packages over ipfs
with a special `XferCommand`.

###### _archlinux:_ reflector

reimplement reflector in go because NIH

#### _k8s:_ Kubernetes

###### _k8s:_

wireguard tunnel for cross cluster connectivity

daemonsets and ip rules

### _Ops_ time and patience needed

Running things other people have written.
Usually just need to learn the model / write config files.

#### _cicd:_ Continuos Integration pipeline

###### _cicd:_ build on push

remove limitation on tagged commits

###### _cicd:_ module proxy for go

run a private module proxy in CI cluster as a local cache

###### _cicd:_ container registry

run a private container registry in CI cluster as a cache for kaniko

###### _cicd:_ deployment for k8s services

as title

#### _archlinux:_ Arch Linux on my laptop

###### _archlinux:_ full disk encryption

add full disk encryption to XPS

### _Write_ wordsmithing wanted

#### _like:_ the [list](/like/) of things i like

###### _like:_ go modules

- [auth](https://github.com/avelino/awesome-go#authentication-and-oauth)
- [cli](https://github.com/avelino/awesome-go#standard-cli)
- [config](https://github.com/avelino/awesome-go#configuration)
- [data](https://github.com/avelino/awesome-go#database)
- [tui](https://github.com/avelino/awesome-go#advanced-console-uis)
- [uuid](https://github.com/avelino/awesome-go#uuid)

#### _blog:_ the notes in my we[b log](/blog/)

###### _blog:_ tekton pipelines

write experience report

###### _blog:_ SNE courses review / learned

experiences / things learned

was it worth it?
