---
description: oops
title: reinstall arch v3
---
just after I finished writing / testing my install script,
I `rm -rf`'d my root directory
(`--no-preserve-root` can't help when root is not at root)

find script here: [arch-install](https://github.com/seankhliao/arch-install)

### usage

1. setup network connectivity
2. setup mount points
   - `mount /dev/nvme0n1p3 /mnt`
   - `mkdir -p /mnt/{efi,boot}`
   - `mount /dev/nvme0n1p1 /mnt/efi`
   - `mkdir -p /mnt/efi/EFI/archlinux`
   - `mount --bind /mnt/efi/EFI/archlinux /mnt/boot`
3. run script
4. `passwd root`
5. `passwd arccy`
6. manually inspect and fix `/etc/fstab`
   - edit to fix bind mount
   - `mkinitcpio`
   - create bootloader entry (or install bootloader)
7. reboot and login
8. setup networking again
9. run `user.sh`
10. setup ssh keys
    - `ssh-keygen -t ed25519`

### notes

#### 30-network.conf

```
[Match]
Name=*

[Network]
DHCP=yes
DNS=8.8.8.8
DNS=8.8.4.4
DNS=2001:4860:4860::8888
DNS=2001:4860:4860::8844
IPForward=kernel
```

#### loader/loader.conf

```
default arch
timeout 0
console-mode max
```

#### loader/entries/arch.conf

```
title 	Arch Linux
linux 	/EFI/archlinux/vmlinuz-linux
initrd	/EFI/archlinux/intel-ucode.img
initrd  /EFI/archlinux/initramfs-linux.img
options root=UUID=...................... quiet rw
```

#### makepkg.conf

```
-march=native
-mtune=native
OPTIONS=(!strip !zipman)
INTEGRITY_CHECK=(sha256)
PKGEXT='.pkg.tar'
SRCEXT='.src.tar'
```