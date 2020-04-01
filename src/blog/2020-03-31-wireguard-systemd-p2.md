---
description: use systemd to manage wireguard
title: wireguard systemd part 2
---

### Goals

- setup a dual stack "server-client vpn" through wireguard using just _systemd-networkd_
- use _public addresses_

#### server

setup interface, give client a /32 ipv4 and a /64 ipv6

```ini
# 40-wg0.netdev
[NetDev]
Name = wg0
Kind = wireguard

[WireGuard]
PrivateKey = SERVER_PRIVATE_KEY
FirewallMark = 1234
ListenPort = 51820

[WireGuardPeer]
PublicKey = CLIENT_PUBLIC_KEY
AllowedIPs = CLIENT_SUBNET_ADDRESS/32,CLIENT_SUBNET_PREFIX::/64
```

setup networking, give server address in shared ranges

- `IPForward` is important and must be set as well as in `/etc/sysctl.d/`

```ini
# 41-wg0.network
[Match]
Name = wg0

[Network]
IPForward = yes
Address = SERVER_SUBNET_ADDRESS/28
Address = SERVER_SUBNET_PREFIX::1/60
```

#### client

##### combined config

setup interface, endpoints are only evaluated once at startup though...

```ini
[NetDev]
Name = wg0
Kind = wireguard

[WireGuard]
PrivateKey = CLIENT_PRIVATE_KEY
FirewallMark = 1234

[WireGuardPeer]
PublicKey = SERVER_PUBLIC_KEY
AllowedIPs = 0.0.0.0/0,::/0
Endpoint = SERVER_PUBLIC_ADDRESS:51820
```

setup networking

```ini
[Match]
Name = wg0

[Network]
Address = CLIENT_SUBNET_ADDRESS/28

[Route]
Destination = 0.0.0.0/0
Table = 2468

[Route]
Destination = ::/0
Table = 2468

[RoutingPolicyRule]
Family = both
InvertRule = true
FirewallMark = 1234
Table = 2468

[RoutingPolicyRule]
Family = both
Table = main
SuppressPrefixLength = 0
```

##### separate ipv4 / ipv6

ex you only need an ipv6 address

- should probably use separate key / peers

###### ipv4

interface

```ini
[NetDev]
Name = wg4
Kind = wireguard

[WireGuard]
PrivateKey = CLIENT_PRIVATE_KEY
FirewallMark = 1234

[WireGuardPeer]
PublicKey = SERVER_PUBLIC_KEY
AllowedIPs = 0.0.0.0/0
Endpoint = SERVER_PUBLIC_ADDRESS:51820
```

networking

```ini
[Match]
Name = wg4

[Network]
Address = CLIENT_SUBNET_ADDRESS/28

[Route]
Destination = 0.0.0.0/0
Table = 2468

[RoutingPolicyRule]
Family = ipv4
InvertRule = true
FirewallMark = 1234
Table = 2468

[RoutingPolicyRule]
Family = ipv4
Table = main
SuppressPrefixLength = 0
```

###### ipv6

interface

```ini
[NetDev]
Name = wg6
Kind = wireguard

[WireGuard]
PrivateKey = CLIENT_PRIVATE_KEY
FirewallMark = 1234

[WireGuardPeer]
PublicKey = SERVER_PUBLIC_KEY
AllowedIPs = ::/0
Endpoint = SERVER_PUBLIC_ADDRESS:51820
```

networking

```ini
[Match]
Name = wg6

[Network]
Address = CLIENT_SUBNET_PREFIX/64

[Route]
Destination = ::/0
Table = 2468

[RoutingPolicyRule]
Family = ipv6
InvertRule = true
FirewallMark = 1234
Table = 2468

[RoutingPolicyRule]
Family = ipv6
Table = main
SuppressPrefixLength = 0
```
