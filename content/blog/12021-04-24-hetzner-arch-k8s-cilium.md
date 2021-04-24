---
title: hetzner arch k8s cilium
description: yay new server
---

### _new_ server

So new Arch Linux dedicated server from Hetzner.
Time to set it up.

```sh
# rename things
NEWHOST=medea
OLDHOST=$(cat /etc/hostname)
echo $NEWHOST > /etc/hostname
sed -i s/$OLDHOST/$NEWHOST/ /etc/hosts

# change mounts / sysctls for k8s
# remove hetzner settings
sed -i '/ swap /d' /etc/fstab
echo br_netfilter > /etc/modules-load.d/br_netfilter.conf
rm /etc/sysctl.d/*
cat << EOF > /etc/sysctl.d/30-ipforward.conf
net.ipv4.ip_forward=1
net.ipv4.conf.lxc*.rp_filter=0
net.ipv6.conf.default.forwarding=1
net.ipv6.conf.all.forwarding=1
EOF

# cleanup junk
pacman -Rns btrfs-progs gptfdisk haveged xfsprogs wget vim net-tools cronie
pacman -S neovim containerd kubeadm kubelet

# remove arch cni dir overrides
rm /etc/kubernetes/kubelet.env
systemctl enable containerd kubelet

sed -i /#PasswordAuthentication yes/PasswordAuthentication no/ /etc/ssh/sshd_config
sed -i /#PrintLastLog yes/PrintLastLog no/ /etc/ssh/sshd_config
rm /etc/ssh/ssh_host_{dsa,rsa,ecdsa}_key*

reboot

kubeadm init --skip-phases=addon/kube-proxy --pod-network-cidr=10.244.0.0/16 --service-cidr=10.245.0.0/16
```

On local machine (see next section for `cilium.yaml`)

```sh
# setup ssh
cat << EOF > ~/.ssh/config

Host medea
    Hostname 65.21.73.144
    # Hostname 2a01:4f9:3b:4e2f::2
    User root
    IdentityFile ~/.ssh/id_ed25519
EOF

ssh-copy-id -i ~/.ssh/id_ecdsa_sk medea
ssh-copy-id -i ~/.ssh/id_ed25519_sk medea

# get kubeconfig
rsync -P medea:/etc/kubernetes/admin.conf ~/.config/kube/config

# allow workloads on my single node
kubectl taint nodes --all node-role.kubernetes.io/master-

# networking
k apply -f cilium.yaml
```

#### _cilium_

So I want some control over the k8s network settings,
specifically the network ranges used and no kube-proxy.
So using the cilium helm chart,
make sure you read the text of the cilium docs to figure out
all the settings you need to enable.

```sh
helm template cilium cilium/cilium \
    --version 1.9.6 \
    --namespace kube-system \
    --set kubeProxyReplacement=strict \
    --set k8sServiceHost=65.21.73.144 \
    --set k8sServicePort=6443 \
    --set operator.replicas=1 \
    --set hostServices.enabled=true \
    --set endpointRoutes.enabled=true \
    --set ipam.mode=cluster-pool \
    --set ipam.operator.clusterPoolIPv4MaskSize=24 \
    --set ipam.operator.clusterPoolIPv4PodCIDR=10.244.0.0/16 \
    --set tunnel=disabled \
    --set autoDirectNodeRoutes=true \
    --set nativeRoutingCIDR=10.244.0.0/15 \
    > cilium.yaml
```
