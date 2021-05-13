---
title: kubevirt notes
description: notes from testing kubevirt
---

### _kubevirt_

[kubevirt](https://kubevirt.io/) is a kubernetes controller
to manage [libvirt](https://libvirt.org/) virtual machines
as CRDs.
It is tightly integreated with its
[containerized-data-importer](https://github.com/kubevirt/containerized-data-importer)
subproject for importing VMs into PVCs.

tested on v0.40.0.

#### _compatibility_

The 2 big things it needs: network and storage.
Networking should work fine, except when you use [cilium](https://cilium.io/) without kube-proxy,
then you run into issues where [services are not reachable](https://github.com/cilium/cilium/issues/14563).
This manifests itself most visibly as DNS not being available from inside the VM.

On storage, I was using [longhorn](https://longhorn.io/)
(who really should fix thier docs to make it more obvious you're looking at an old version)
and ran into a problem copying between PVCs, tldr, mounting as readonly
is [not yet supported](https://github.com/longhorn/longhorn/issues/2575).
Oh well, host the image somewhere and download it every time I guess.

#### _using_

`kubectl virt console myvm` never actually worked for me,
vnc was much more reliable.

On the images themselves, using [cloud-init](https://cloudinit.readthedocs.io/en/latest/)
is probably the safest option to get started,
such as [Arch Linux VM images](https://gitlab.archlinux.org/archlinux/arch-boxes/-/jobs/21699/artifacts/browse/output)
(choose cloudimg).
And if you never used it: the main entrypoint can be a script or a config file,
identified with `#cloud-init` in place of a shebang.

On making the images, I'm partial to [mkosi](https://github.com/systemd/mkosi)
since creating the image doesn't actually require spinning up a vm.
