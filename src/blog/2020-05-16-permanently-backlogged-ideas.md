---
description: ideas that i will probably never get round to implementing
title: permanently backlogged ideas
---

### _Dev_ implementation required

Ideas that need code to be written,
may or may not see the light of day

##### _go:_ gopher things

- **http middleware handlers:**
  find a nice list of things,
  ex, add headers, add compression,
  add security(csp, cors)...
- **module proxy over IPFS:**
  modify `go` command to retrieve modules over IPFS
  maybe works well with IPFS deduplication?

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

### _Ops_ time and patience needed

Running things other people have written.
Usually just need to learn the model / write config files.

##### _cicd:_ Continuous Integration pipeline

- **build on push/pull request:**
  remove limitation on tagged commits,
  lint, test, ???

##### _archlinux:_ Arch Linux on my laptop

- **full disk encryption:**
  add full disk encryption to XPS

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
  run [gerrit](https://www.gerritcodereview.com/) as private git server
- **code search:**
  run [sourcegraph](https://github.com/sourcegraph/sourcegraph) on private git server
