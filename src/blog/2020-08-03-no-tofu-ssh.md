---
description: eliminate trust on first use for ssh
title: no tofu ssh
---

### _ssh_ trust on first use

ssh into a new machine (or one that has recently changed keys)
and we get the familiar message:

```txt
The authenticity of host 'sshfp.seankhliao.com (34.xx.xx.xx)' can't be established.
ECDSA key fingerprint is SHA256:BQYZU02snggz6PD5ZUfjigqb/ZxswcgnSzHkk8/PVD8.
Are you sure you want to continue connecting (yes/no/[fingerprint])?
```

so what can you do about this?

- [sshfp](#SSHFP)
- [certs](#certificates)

#### _SSHFP_

SSH FingerPrint uses DNSSEC enabled DNS as a secure vector to determine trust.
It assumes that a host compromise doesn't include a DNS compromise,
so whatever is in the DNS records can be trusted as the proper keys.

##### _server_ side setup

generate the SSHFP DNS records

```sh
ssh-keygen -r host.example.com.
```

output

```dns
sshfp IN SSHFP 1 1 70e809125961a456bbd8178cf8e9a0d4addc32fa
sshfp IN SSHFP 1 2 e5b848dcd7ab7fea8e6f8db217ab02bdba93c5b270bb6ece1a707e40f98a842a
sshfp IN SSHFP 2 1 5177759e837130e3f14276e27c2721984dc08039
sshfp IN SSHFP 2 2 98377672b8ee14286e9f962a985923e66acbc8bb03113c4c7fad62ceb5c809b5
sshfp IN SSHFP 3 1 ba4bb33c362ff955d13d71881f0ae19b15f558f7
sshfp IN SSHFP 3 2 050619534dac9e0833e8f0f96547e38a0a9bfd9c6cc1c8274b31e493cfcf543f
sshfp IN SSHFP 4 1 5f31fd8a762b2381e87a718ca8a77495967854d8
sshfp IN SSHFP 4 2 b93e879569b28d1769f8962c479eb9f52b734b3f17ff06bad721b8c0a42221ff
```

##### _DNS_ setup

- point a A/AAAA record to host
- copy the SSHFP records from server side in

##### _client_ side setup

- set `VerifyHostKeyDNS yes` in ssh config
  (global `/etc/ssh/ssh_config`, per user `~/.ssh/config`)

#### _certificates_

The basic idea is the same:
you trust a certificate authority
and by extension you trust everything it signs.

SSH certs/ca does not have a hierarchy,
CA directly signs user/client.

##### _server_ certificates

Servers present certificate, client trusts CA

###### _DNS_ setup

point a A/AAAA record to host

###### _Certicate Authority_ setup

On a trusted computer generate a CA SSH keypair,
example using security key:

```sh
ssh-keygen -t ecdsa-sk -f host-ca
```

output

```txt
Generating public/private ecdsa-sk key pair.
You may need to touch your authenticator to authorize key generation.
Enter PIN for authenticator:
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in host-ca
Your public key has been saved in host-ca.pub
The key fingerprint is:
SHA256:uQlL6dweX5yKwjgPBInCPiHV4tBataXq/v+j/EvW69Q arccy@eevee
The key's randomart image is:
+-[ECDSA-SK 256]--+
| ..o. .          |
|oo+..+           |
|+Bo.o            |
|= oo   . .       |
| o. . + S        |
| ... + + + o .   |
|  . .o+ B o E    |
| .  ooo+.= +     |
|  ...+==*+=      |
+----[SHA256]-----+
```

obtain a host's public key and sign it,
note the use of the host's DNS name for `-n`

```sh
ssh-keygen -h -I instance-1 -s host-ca -V +52w -z 1 -n sshfp.seankhliao.com ssh_host_ed25519_key.pub
```

###### _server_ side setup

copy the generated cert back to host and
add `HostCertificate /etc/ssh/ssh_host_ed25519_key.pub` to sshd config (`/etc/ssh/sshd_config`)

###### _client_ side setup

trust the certificate authority in `known_hosts`:

```txt
@cert-authority sshfp.seankhliao.com sk-ecdsa-sha2-nistp256@openssh.com AAAAInNr...`
```

###### _revoking_ server certs

revoke an entire CA by changing `@cert-authority` to `@revoke`

not sure if possible to revoke single hosts

##### _client_ certificates

Clients present certificates, server trusts CA

no more pesky ssh-copy-id dance,
setup a host once and login for any trusted user

###### _Certicate Authority_ setup

On a trusted computer generate a CA SSH keypair,
example using security key:

```sh
ssh-keygen -t ecdsa-sk -f client-ca
```

output

```txt
Generating public/private ecdsa-sk key pair.
You may need to touch your authenticator to authorize key generation.
Enter PIN for authenticator:
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in client-ca
Your public key has been saved in client-ca.pub
The key fingerprint is:
SHA256:6TGaGBZryUtDzglSdtMh27j2RjqJpgmD0KYVal6h2Sc arccy@eevee
The key's randomart image is:
+-[ECDSA-SK 256]--+
|  o +...         |
| o . *.          |
|. o.= .          |
| ++B.*   .       |
|o++E/.. S        |
|=+.*oX + o       |
|=.o * = .        |
|.=   o           |
|o                |
+----[SHA256]-----+
```

obtain client public key and sign it with

```sh
ssh-keygen -I seankhliao -s client-ca -z 2 -V +52w -n arccy id_client.pub
```

###### _client_ side setup

copy the generated cert back to client
and specify using it with `CertificateFile ~/.ssh/id_client-cert.pub`
and the private ket file with `IdentityFile ~/.ssh/id_client`

###### _server_ side setup

copy the CA pubkey over
and trust it in `sshd_config` with `TrustedUserCAKeys /etc/ssh/client-ca.pub`

###### _revoking_ client certs

remove entire CA from `TrustedUserCAKeys`

not sure if possible to revoke single clients
