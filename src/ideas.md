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

##### _cicd:_ Continuos Integration Pipeline

- **tekton webhook interceptor:**
  cel is nice but a bit too limiting,
  implement custom interceptor with plugin style
- **get latest version:**
  retrieve the latest version of thrid party X
  from a variety of sources (github, arch, other).
  update / tag repo and push for ci.
  partially done for `ci-*` projects on github (with actions)

##### _go:_ go core things

- **module proxy over IPFS:**
  modify `go` command to retrieve modules over IPFS
  maybe works well with IPFS deduplication?

##### _web:_ webdev things

- **webauthn:**
  try it

##### _site:_ this website

- **serve over AMP:**
  this site doesn't need fancy things,
  generate and serve over AMP with subdomain & use with cloudflare
- **serve with WebPackaging / SXG:**
  serve signed / shareable bundle for site / page,
  needs proper SXG cert: expensive
- **serve over IPFS:**
  use DNSLink and serve this site over IPFS
- **serve over Tor alt-svc:**
  apparently it's not difficult?

##### _archlinux:_ Arch Linux Infra

- **mirror/repo over IPFS:**
  Arch should be flexible enough to handle getting packages over ipfs
  with a special `XferCommand`.
- **reflector:**
  reimplement reflector in go because NIH

##### _k8s:_ Kubernetes

- **cross cluster tunnel:**
  wireguard tunnel for cross cluster connectivity,
  with daemonsets and ip rules
- **terminate tls at pods:**
  l4 ingress and terminate at pods instead,
  requires a better ca setup for internal services,
  or split horizon dns
- **module proxy for go:**
  run [goproxy](https://github.com/goproxy/goproxy) in CI cluster as a local cache
- **container registry:**
  run [harbor](https://github.com/goharbor/harbor) in CI cluster as a cache for kaniko
- **git server:**
  run private git server:
  gerrit? gitea ? plain git?

### _Ops_ time and patience needed

Running things other people have written.
Usually just need to learn the model / write config files.

##### _cicd:_ Continuos Integration pipeline

- **build on push:**
  remove limitation on tagged commits
- **deployment for k8s services:**
  as title

##### _archlinux:_ Arch Linux on my laptop

- **full disk encryption:**
  add full disk encryption to XPS

### _Write_ wordsmithing wanted

##### _blog:_ the notes in my we[b log](/blog/)

- **tekton pipelines:**
  write experience report
- **SNE courses review / learned:**
  experiences / things learned,
  was it worth it?
