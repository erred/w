---
description: so many levels of abstraction for container runtimes
title: container runtime
---

### _Container_ Runtimes

So many layers of abstraction.
But tl;dr:

K8s / docker - containerd / cri-o - runtime

#### _Tools_

Tools end users use to run containers

##### _k8s_ Kubernetes

[k8s][k8s] is an orchestrator for
scheduling/running containers across multiples nodes (machines).

Supports different runtimes
declared in [`RuntimeClass`][runtimeclass] in `node.k8s.io/v1beta1`.
Runtimes can be selected through `runtimeClassName` in a PodSpec.

The interface for pluggable runtimes to support is
[CRI][#cri-container-runtime-interface].

##### _Docker_

[docker][docker] is a engine/daemon/cli for running containers
on a single machine.

Supports different runtimes declared in `runtimes` in `daemon.json`.
Runtimes can be selected through `--runtime=x` or `default-runtime` in `daemon.json`.

The interface for pluggable runtimes to support is
[OCI][oci-open-container-initiative-runtime-specification].

##### _Others_

- [podman][podman] similar to docker, supports OCI
- [LXC][lxc] / LXD (doesn't support OCI, focus on supporting entire OS)
- [systemd-nspawn][nspawn] supports OCI, low user traction
- [Mesos][mesos]
- [OpenVZ][openvz]
- [rkt][rkt] (deprecated)

#### _Interfaces_

##### _CRI_ Container Runtime Interface

k8s standard for high level runtimes.
Not sure of any practical difference between the 2 main implementations.

###### _containerd_

Broken out of docker,
[`containerd`][containerd] is the high level daemon
responsible for the container lifecycle.

CNCF Graduated, supports low level OCI runtimes
through a [shim][containerd-shim]

###### _CRI-O_

Developed by RedHat,
[`cri-o`][cri-o] / `crio` is a lightweight alternative
to running containers.

CNCF Incubating, supports OCI runtimes.

###### _other_ runtimes

[docker][docker] through [`dockershim`][dockershim],
adds extra layer between k8s and containerd,
doesn't allow for different runtime classes.
Likely to be deprecated in future?

[frakti][frakti] is a runtime for VMs.

##### _OCI_ Open Container Initiative Runtime Specification

[OCI Runtime Spec][oci] is the standard for low level container runtimes,
detailing how to run a "filesystem bundle".
Not to be confusted with [OCI Image Spec][ociimg].

###### _runc_

[runc][runc] is the default/reference implementation
for low level container runtimes.

###### _gVisor_

[gVisor][gvisor] is a process kernel by Google in Go.

Use as a standalone low level runtime as [runsc][runsc]
or with containerd [containerd-shim-runsc-v1][containerd-runsc]

###### _Kata_ Containers

Lightweight VMs?

Use as a standalone low level runtime as `kata-runtime`
or with containerd [io.containerd.kata.v2][kata-containerd]

###### _Nabla_ Containers

Not sure what it does to be "more isolated".

Use as a standalone low level runtime as [runnc][runnc].

###### _firecracker_

MicroVM by Amazon in Rust.

Use with Kata [kata-fc][katafc]
or with containerd [firecracker-containderd][firecracker-containderd].

###### _others_

- [bwrap-oci][bwrap] (unmaintained)
- [clearcontainers][clear] (deprecated in favour of kata containers)
- [crio-lxc][criolxc] wrapper for CRI-O to control LXC
- [crun][crun] written in C
- [Railcar][railcar] (deprecated)
- [rkt][rkt] (deprecated)
- runlxc currently proprietary runtime by Alibaba
- [runV][runv] (unmaintained)

[bwrap]: https://github.com/projectatomic/bwrap-oci
[clear]: https://github.com/clearcontainers
[containerd]: https://containerd.io/
[containerd-runsc]: https://gvisor.dev/docs/user_guide/containerd/quick_start/
[containerd-shim]: https://github.com/containerd/containerd/tree/master/runtime/v2
[cri-o]: https://cri-o.io/
[criolxc]: https://github.com/lxc/crio-lxc
[crun]: https://github.com/containers/crun
[docker]: https://www.docker.com/
[dockershim]: https://github.com/kubernetes/kubernetes/tree/master/pkg/kubelet/dockershim
[firecracker]: https://firecracker-microvm.github.io/
[firecracke-containderd]: https://github.com/firecracker-microvm/firecracker-containerd
[frakti]: https://github.com/kubernetes/frakti
[gvisor]: https://gvisor.dev/
[k8s]: https://kubernetes.io/
[kata]: https://katacontainers.io/
[kata-containerd]: https://github.com/kata-containers/documentation/blob/master/how-to/how-to-use-k8s-with-cri-containerd-and-kata.md
[katafc]: https://github.com/kata-containers/documentation/wiki/Initial-release-of-Kata-Containers-with-Firecracker-support
[lxc]: https://linuxcontainers.org/
[mesos]: http://mesos.apache.org/
[nabla]: https://nabla-containers.github.io/
[nspawn]: https://www.freedesktop.org/software/systemd/man/systemd-nspawn.html
[oci]: https://github.com/opencontainers/runtime-spec
[ociimg]: https://github.com/opencontainers/image-spec
[openvz]: https://openvz.org/
[podman]: https://podman.io/
[railcar]: https://github.com/oracle/railcar
[rkt]: https://github.com/rkt/rkt
[runc]: https://github.com/opencontainers/runc
[runnc]: https://github.com/nabla-containers/runnc
[runsc]: https://gvisor.dev/docs/user_guide/quick_start/docker/
[runtimeclass]: https://kubernetes.io/docs/concepts/containers/runtime-class/
[runv]: https://github.com/hyperhq/runv
