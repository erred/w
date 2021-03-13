---
description: starting a kubernetes cluster
title: k8s cluster setup
---

### _kubernetes_

_tldr:_ `kind` for local testing, `k3s` for dead easy cluster setup, `kubeadm` for a customized install

all are tested with a Arch Linux host created with:

```sh
gcloud compute instances create k0 \
    --image-project arch-linux-gce --image-family arch \
    --machine-type e2-medium
```

#### _single_ machine

##### _kind_

- use `extraPortMappings` to map host ports onto a node, then use an ingresscontroller with hostports

```sh
sudo pacman -Syu docker kubectl
sudo systemctl enable --now docker
sudo usermod -aG docker $(whoami)
newgrp docker

sudo curl -Lo /usr/local/bin/kind https://kind.sigs.k8s.io/dl/v0.9.0/kind-linux-amd64
sudo chmod +x /usr/local/bin/kind
kind create cluster

kubectl get --all-namspaces all
```

##### _minikube_

[minikube](https://minikube.sigs.k8s.io/docs/)

- separate command to tunnel into LoadBalancer services

```sh
sudo pacman -Syu docker kubectl
sudo systemctl enable --now docker
sudo usermod -aG docker $(whoami)
newgrp docker

sudo curl -Lo /usr/local/bin/minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo chmod +x /usr/local/bin/minikube
minikube start

kubectl get --all-namspaces all
```

##### _microk8s_

This needs snap... too much of a PITA to install

#### _multi_ machine

##### _k3s_

[k3s](https://k3s.io/)

- starts traefik with hostports 80, 443, 8080
- Klipper LB (node hostports) for Service LoadBalancer
- installs itself as a systemd service

```sh
sudo pacman -Syu which

curl -sfL https://get.k3s.io | sh -

kubectl get --all-namspaces all
```

##### _kubeadm_

- provide your own stuff, this is "just" responsible for half the config
- using the containerd runtime, `--cni-bin-dir` isn't supported, either:
  - configure containerd to match the kubelet args [`/etc/containerd/config.toml`](https://github.com/containerd/cri/blob/master/docs/config.md)
  - make sure kubelet uses the default `/opt/cni/bin`

```sh
sudo pacman -Syu containerd kubeadm kubelet kubectl
sudo rm /etc/sysctl.d/60-gce-network-security.conf
sudo sysctl --system
sudo modprobe br_netfilter
sudo mkdir -p /etc/containerd
cat << EOF | sudo tee /etc/containerd/config.toml
version = 2
[plugins."io.containerd.grpc.v1.cri".cni]
  bin_dir = "/usr/lib/cni"
EOF
sudo systemctl enable --now containerd
sudo systemctl enable kubelet

sudo kubeadm init --pod-network-cidr=10.244.0.0/16
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

sudo kubectl --kubeconfig /etc/kubernetes/admin.conf get --all-namespaces all
```

###### _cilium_

- [system requirements](https://docs.cilium.io/en/v1.9/operations/system_requirements/) is really important, ex that note about systemd and `rp_filter`

```sh
sudo pacman -Syu containerd kubeadm kubelet kubectl
sudo rm /etc/sysctl.d/60-gce-network-security.conf
echo 'net.ipv4.conf.lxc*.rp_filter = 0' | sudo tee /etc/sysctl.d/99-override_cilium_rp_filter.conf
sudo systemctl restart systemd-sysctl
sudo modprobe br_netfilter
sudo mkdir -p /etc/containerd
echo 'KUBELET_ARGS=' | sudo tee /etc/kubernetes/kubelet.env
sudo systemctl enable --now containerd
sudo systemctl enable kubelet

sudo kubeadm init
kubectl create -f https://raw.githubusercontent.com/cilium/cilium/v1.9/install/kubernetes/quick-install.yaml

sudo kubectl --kubeconfig /etc/kubernetes/admin.conf get --all-namespaces all
```

##### _k0s_

[k0s](https://k0sproject.io/)

- doesn't do much for networking / services

```sh
sudo pacman -Syu kubectl
echo 127.0.0.1 localhost | sudo tee -a /etc/hosts

sudo curl -Lo /usr/local/bin/k0s https://github.com/k0sproject/k0s/releases/download/v0.8.1/k0s-v0.8.1-amd64
sudo chmod +x /usr/local/bin/k0s
sudo k0s server

sudo kubectl --kubeconfig /var/lib/k0s/pki/admin.conf get --all-namespaces all
```

##### _kops_

[kOps](https://kops.sigs.k8s.io/)

This is for creating your self managed kubernetes clusters in cloud environments,
not exactly what I want to testing today.
Short setup though.

#### _kubespray_

[kubespray](https://kubespray.io/#/)

Similar to kOps but uses ansible and slightly more generic?
Very long setup.
