---
title: terraform single server
description: using terraform for a bare metal server
---

### _terraform_

[terraform](https://www.terraform.io/) rightly describes itself as a provisioning tool,
and not a config management tool.
So this idea of using terraform to manage the config for a server didn't really work out.

#### _plan_

use `provisioner "file"`s on `null_resource`s to transfer the config files
necessary to run [nomad](https://www.nomadproject.io/).

#### _results_

You have to do it right the first time,
otherwise it's a manual `taint` of the resource and recreating it.
Also it doesn't create the directories for a specified path,
necessitating `provisioner "remote-exec"` just to create the dirs.

I should (have to?) reboot the server after everything is done,
but terraform doesn't make doing that + exiting properly / validating easy.
So I copped out and just did that part by hand.

#### _files_

##### 00_certs.tf

This file creates a self signed CA
and a set of client certs output to a local dir
so local `nomad` cli commands can authenticate.

```terraform
##################################################
# ca certs
##################################################
resource "tls_private_key" "ca" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_self_signed_cert" "ca" {
  key_algorithm     = tls_private_key.ca.algorithm
  private_key_pem   = tls_private_key.ca.private_key_pem
  is_ca_certificate = true

  validity_period_hours = 24 * 365 * 10

  allowed_uses = [
    "cert_signing",
  ]

  subject {
    common_name  = "medea.seankhliao.com"
    organization = "medea / seankhliao"
  }
}
resource "local_file" "ca_cert" {
  content  = tls_self_signed_cert.ca.cert_pem
  filename = "${path.root}/out/ca.crt"
}


##################################################
# client cert
##################################################
resource "tls_private_key" "nomad_client" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}
resource "local_file" "nomad_client_key" {
  content  = tls_private_key.nomad_client.private_key_pem
  filename = "${path.root}/out/client.key"
}

resource "tls_cert_request" "nomad_client" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.nomad_client.private_key_pem

  subject {
    common_name  = "eevee.seankhliao.com"
    organization = "medea / seankhliao"
  }
}
resource "tls_locally_signed_cert" "nomad_client" {
  ca_key_algorithm   = tls_private_key.ca.algorithm
  ca_private_key_pem = tls_private_key.ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem
  cert_request_pem   = tls_cert_request.nomad_client.cert_request_pem

  validity_period_hours = 24 * 365

  allowed_uses = [
    "digital_signature",
    "key_encipherment",
    "client_auth",
  ]
}
resource "local_file" "nomad_client_cert" {
  content  = tls_locally_signed_cert.nomad_client.cert_pem
  filename = "${path.root}/out/client.crt"
}
```

##### 01_medea.tf

[medea](https://your-throne.fandom.com/wiki/Medea_Solon)
is the name of my server.

This file manages the basic setup i need from a fresh image.

```terraform
resource "null_resource" "medea" {
  connection {
    host        = "medea.seankhliao.com"
    private_key = file(pathexpand("~/.ssh/id_ed25519"))
    agent       = false
  }

  provisioner "remote-exec" {
    inline = [
      "rm /etc/sysctl.d/* || true",
      "rm /etc/ssh/ssh_host_{dsa,rsa,ecdsa}_key* || true",
      "pacman -Rns --noconfirm btrfs-progs gptfdisk haveged xfsprogs wget vim net-tools cronie || true",
      "pacman -Syu --noconfirm neovim zip unzip",
      "systemctl enable --now systemd-timesyncd",
    ]
  }
  provisioner "file" {
    destination = "/root/.ssh/authorized_keys"
    content     = <<-EOT
      ${file(pathexpand("~/.ssh/id_ed25519.pub"))}
      ${file(pathexpand("~/.ssh/id_ed25519_sk.pub"))}
      ${file(pathexpand("~/.ssh/id_ecdsa_sk.pub"))}
    EOT
  }
  provisioner "file" {
    destination = "/etc/sysctl.d/30-ipforward.conf"
    content     = <<-EOT
      net.ipv4.ip_forward=1
      net.ipv4.conf.lxc*.rp_filter=0
      net.ipv6.conf.default.forwarding=1
      net.ipv6.conf.all.forwarding=1
    EOT
  }
  provisioner "file" {
    destination = "/etc/modules-load.d/br_netfilter.conf"
    content     = "br_netfilter"
  }
  provisioner "file" {
    destination = "/etc/systemd/network/40-wg0.netdev"
    content     = <<-EOT
      # WireGuard

      [NetDev]
      Name = wg0
      Kind = wireguard

      [WireGuard]
      PrivateKey = xxxx
      # PublicKey = lombY0b15giOmoM9t0xBi+UgVkZDoOKDaEV9+ONwH1U=
      ListenPort = 51820

      # eevee
      [WireGuardPeer]
      PublicKey = YvSLDXl3NX1ySTX2C8D72+fCVBcqSs+fmAX3uySCDAQ=
      AllowedIPs = 192.168.100.13/32

      # pixel 3
      [WireGuardPeer]
      PublicKey = 3xpGlOORQb9/yg545KX+odrup3YaslxO9ie+ztJ3Y3E=
      AllowedIPs = 192.168.100.3/32

      # arch
      [WireGuardPeer]
      PublicKey = Lr17jGvc7uwjn9LNRR+IkCkjuP8nkHTOMHbVV+onMn0=
      AllowedIPs = 192.168.100.1/24
      Endpoint = yyyy
    EOT
  }
  provisioner "file" {
    destination = "/etc/systemd/network/41-wg0.network"
    content     = <<-EOT
      [Match]
      Name = wg0

      [Network]
      Address = 192.168.100.2/24
    EOT
  }
}
```

##### 02_nomad.tf

This file provisions the server(+client) certs for nomad
(it needs client to talk to itself),
downloads nomad and runs it with `systemd`.

```terraform
##################################################
# nomad server cert
##################################################
resource "tls_private_key" "nomad_server" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_cert_request" "nomad_server" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.nomad_server.private_key_pem

  dns_names = [
    "localhost",
    "server.global.nomad",
    "medea.seankhliao.com",
  ]

  ip_addresses = [
    "127.0.0.1",
    "192.168.100.2",
    "65.21.73.144",
  ]

  subject {
    common_name  = "medea.seankhliao.com"
    organization = "medea / seankhliao"
  }
}
resource "tls_locally_signed_cert" "nomad_server" {
  ca_key_algorithm   = tls_private_key.ca.algorithm
  ca_private_key_pem = tls_private_key.ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem
  cert_request_pem   = tls_cert_request.nomad_server.cert_request_pem

  validity_period_hours = 24 * 365

  allowed_uses = [
    "digital_signature",
    "key_encipherment",
    "server_auth",
    "client_auth",
  ]
}

##################################################
# nomad
##################################################
resource "null_resource" "medea_nomad" {
  depends_on = [
    null_resource.medea,
  ]

  connection {
    host        = "medea.seankhliao.com"
    private_key = file(pathexpand("~/.ssh/id_ed25519"))
    agent       = false
  }

  provisioner "file" {
    destination = "/etc/systemd/system/nomad.service"
    content     = <<-EOT
      [Unit]
      Description=nomad server
      Documentation=https://www.nomadproject.io/docs/agent/
      Requires=network-online.target
      After=network-online.target

      [Service]
      Restart=on-failure
      ExecStart=/usr/local/bin/nomad agent -config /etc/nomad/server.hcl
      ExecReload=/usr/bin/kill -HUP $MAINPID
      KillSignal=SIGINT

      [Install]
      WantedBy=multi-user.target
    EOT
  }
  provisioner "remote-exec" {
    inline = [
      "curl -Lo /tmp/nomad.zip https://releases.hashicorp.com/nomad/1.1.0/nomad_1.1.0_linux_amd64.zip",
      "unzip -o /tmp/nomad.zip nomad -d /usr/local/bin",
      "systemctl daemon-reload",
      "systemctl enable nomad",
      "rm -rf /etc/nomad || true",
      "mkdir -p /etc/nomad",
    ]
  }
  provisioner "file" {
    destination = "/etc/nomad/ca.crt"
    content     = tls_self_signed_cert.ca.cert_pem
  }
  provisioner "file" {
    destination = "/etc/nomad/server.key"
    content     = tls_private_key.nomad_server.private_key_pem
  }
  provisioner "file" {
    destination = "/etc/nomad/server.crt"
    content     = tls_locally_signed_cert.nomad_server.cert_pem
  }
  provisioner "file" {
    destination = "/etc/nomad/server.hcl"
    content     = <<-EOT
      bind_addr = "0.0.0.0"
      data_dir  = "/var/lib/nomad"

      tls {
        http = true
        rpc  = true

        ca_file   = "/etc/nomad/ca.crt"
        cert_file = "/etc/nomad/server.crt"
        key_file  = "/etc/nomad/server.key"

        verify_server_hostname = true
        verify_https_client    = true
      }

      leave_on_interrupt = true
      leave_on_terminate = true

      disable_update_check = true

      server {
        enabled = true
        bootstrap_expect = 1
      }

      client {
        enabled = true
      }

      log_level = "INFO"
    EOT
  }
}
```

##### 03_consul.tf

Similar, but for consul.
The other option is to run it inside nomad as a system job,
which may actually work better, since this is janky.
Once nomad has talked to consul, it doesn't like being disconnected.

```terraform
##################################################
# consul server cert
##################################################
resource "tls_private_key" "consul_server" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_cert_request" "consul_server" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.consul_server.private_key_pem

  dns_names = [
    "localhost",
    "server.dc1.consul",
    "medea.seankhliao.com",
  ]

  ip_addresses = [
    "127.0.0.1",
    "192.168.100.2",
    "65.21.73.144",
  ]

  subject {
    common_name  = "medea.seankhliao.com"
    organization = "medea / seankhliao"
  }
}
resource "tls_locally_signed_cert" "consul_server" {
  ca_key_algorithm   = tls_private_key.ca.algorithm
  ca_private_key_pem = tls_private_key.ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.ca.cert_pem
  cert_request_pem   = tls_cert_request.consul_server.cert_request_pem

  validity_period_hours = 24 * 365

  allowed_uses = [
    "digital_signature",
    "key_encipherment",
    "server_auth",
    "client_auth",
  ]
}

##################################################
# consul
##################################################
resource "null_resource" "medea_consul" {
  depends_on = [
    null_resource.medea,
  ]

  connection {
    host        = "medea.seankhliao.com"
    private_key = file(pathexpand("~/.ssh/id_ed25519"))
    agent       = false
  }

  provisioner "file" {
    destination = "/etc/systemd/system/consul.service"
    content     = <<-EOT
      [Unit]
      Description=Consul Agent
      Documentation=https://consul.io/docs/
      Requires=network-online.target
      After=network-online.target

      [Service]
      Restart=on-failure
      ExecStart=/usr/local/bin/consul agent -config-file /etc/consul/server.hcl
      ExecReload=/usr/bin/kill -HUP $MAINPID
      KillSignal=SIGINT

      [Install]
      WantedBy=multi-user.target
    EOT
  }
  provisioner "remote-exec" {
    inline = [
      "curl -Lo /tmp/consul.zip https://releases.hashicorp.com/consul/1.9.5/consul_1.9.5_linux_amd64.zip",
      "unzip -o /tmp/consul.zip consul -d /usr/local/bin",
      "systemctl daemon-reload",
      "systemctl enable consul",
      "rm -rf /etc/consul || true",
      "mkdir -p /etc/consul",
    ]
  }
  provisioner "file" {
    destination = "/etc/consul/ca.crt"
    content     = tls_self_signed_cert.ca.cert_pem
  }
  provisioner "file" {
    destination = "/etc/consul/server.key"
    content     = tls_private_key.consul_server.private_key_pem
  }
  provisioner "file" {
    destination = "/etc/consul/server.crt"
    content     = tls_locally_signed_cert.consul_server.cert_pem
  }
  provisioner "file" {
    destination = "/etc/consul/server.hcl"
    content     = <<-EOT
      bind_addr = "0.0.0.0"
      data_dir  = "/var/lib/consul"
      ports {
        http = -1
        https = 8501
        grpc = 8502
      }

      ca_file   = "/etc/consul/ca.crt"
      cert_file = "/etc/consul/server.crt"
      key_file  = "/etc/consul/server.key"

      verify_incoming = true
      verify_server_hostname = true

      leave_on_terminate = true

      disable_update_check = true

      server = true
      bootstrap_expect = 1

      client_addr = "0.0.0.0"

      ui_config {
        enabled = true
      }

      log_level = "INFO"
    EOT
  }
}
```
