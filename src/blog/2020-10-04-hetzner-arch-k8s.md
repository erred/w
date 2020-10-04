---
description: setup for hetzner archlinux boxes to run k8s
title: hetzner arch k8s
---

### _hetzner_

So [Hetzner](https://www.hetzner.com/) sells a bunch of computing stuff.
I went for their dedicated servers.

_note_: slightly cheaper for compute+storage compared to buying committed use on GCP.
free networking is _+++_.

Goals: setup kubernetes so I have a mini cloud to play with.

#### _Arch_ Linux

##### _root_

Hetzner has the option to install "minimal" Arch Linux as the OS,
imo, not minimal enough.

Step 1: clean out unneeded packages and install a few necessary ones

_note_: `neovim` vs `vim` is just preference,
`zsh` is also preference,
`binutils` and `fakeroot` are needed for installing stuff from AUR,
`rsync` might already be installed,
`sudo` because I don't want to be root all the time,
`fail2ban` is not strictly necessary but helps keep logs clean.

```sh
pacman -Rns btrfs-progs gptfdisk haveged xfsprogs wget vim net-tools cronie
pacman -Syu sudo zsh git fail2ban neovim binutils fakeroot kubectl rsync docker

systemctl enable --now fail2ban
systemctl enable --now docker
```

Step 2: ensure needed forwarding conf:

```sh
echo br_netfilter | sudo tee /etc/modules-load.d/br_netfilter.conf
cat << EOF | sudo tee /etc/sysctl.d/30-ipforward.conf
net.ipv4.ip_forward=1
net.ipv6.conf.default.forwarding=1
net.ipv6.conf.all.forwarding=1
EOF
```

Step 3: manually edit some files

```sh
# edit /etc/hostname

# edit /etc/hosts

# edit /etc/fstab
#   disable swap

# edit /etc/ssh/sshd_config
#   PermitRootLogin no
#   PasswordAuthentication no
#   PrintLastLog no
```

Step 5: create user and add to sudo group

_note_: the defaults are apparently different than usual like no create home?

_note_: docker group is root-equivalent

```sh
echo '%sudo ALL=(ALL) ALL' > sudoers.d/sudo
groupadd sudo
useradd -G sudo,docker -s /bin/zsh -m arccy
passwd arccy
```

from own machine copy ssh keys for login

```sh
ssh-copy-id -i ~/.ssh/id_ecdsa_sk hetz
ssh-copy-id -i ~/.ssh/id_ed25519_sk hetz
ssh-copy-id -i ~/.ssh/id_ed25519 hetz
```

Step 6: reboot for some of the changes to take effect

```sh
systemctl reboot
```

##### _user_

actions will be run a user, so login

Step 7: copy configs for familiarity

```sh
git clone https://github.com/seankhliao/config .config
ln -s .config/zsh/zshenv .zshenv
mkdir -p /home/arccy/data/xdg
```

Step 8: _optional_ create and add ssh keys for github,
else disable url rewrites in gitconfig

```sh
ssh-keygen -t ed25519
# add key to github
```

Step 9: install stuff from AUR

_note_: not strictly necessary, maybe better to install out of band? without this, no need for `yay` either

```sh
curl -LO https://github.com/Jguer/yay/releases/download/v10.0.4/yay_10.0.4_x86_64.tar.gz
tar xf yay_10.0.4_x86_64.tar.gz
cd yay_10.0.4_x86_64
chmod +x yay
./yay -Syu yay-bin kind-bin
cd
rm -rf yay* .bash* .viminfo
```

Step 10: create config and cluster

This uses [kind](https://kind.sigs.k8s.io/) to run a cluster,
originally the plan was a kubeadm installed cluster + cilium + istio
but it wasn't working out too well
(dunno why, it worked last time).

Disadvantage is hostPorts need to be mapped ahead of time.

_note_: alternate port mapping strategy would map into the `NodePort`range
so services can be targeted and hostPort+DaemonSet unnecessary.

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: 135.181.76.154
nodes:
  - role: control-plane
    extraPortMappings:
      # http3
      - containerPort: 443
        hostPort: 443
        protocol: udp

        # tls, https
      - containerPort: 443
        hostPort: 443
        protocol: tcp

        # http
      - containerPort: 80
        hostPort: 80
        protocol: tcp

        # ipfs swarm
      - containerPort: 4001
        hostPort: 4001
        protocol: udp
      - containerPort: 4001
        hostPort: 4001
        protocol: tcp

        # wireguard
      - containerPort: 51820
        hostPort: 51820
        protocol: udp
```
