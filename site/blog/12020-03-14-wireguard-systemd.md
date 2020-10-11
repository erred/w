---
description: use systemd to manage wireguard
title: wireguard systemd
---

### Goals

- setup a default VPN through wireguard using just _systemd-networkd_

##### _additional_ reading

- ArchWiki [WireGuard/Systemd](https://wiki.archlinux.org/index.php/WireGuard#Using_systemd-networkd)
- WireGuard [WireGuard/Improved Rule](https://www.wireguard.com/netns/#improved-rule-based-routing)

#### How

Use the improved rule-based solution similar to wg-quick

##### 40-wireguard.netdev

creating the network device

- **WireGuard**: sets _mark 1234_ on outgoing packets

```ini
[NetDev]
Name = wg0
Kind = wireguard

[WireGuard]
PrivateKey = CLIENT_PRIVATE_KEY
FirewallMark = 1234

[WireGuardPeer]
PublicKey = SERVER_PUBLIC_KEY
AllowedIPs = 0.0.0.0/0
Endpoint = SERVER_ENDPOINT:51820
```

##### 40-wireguard.network

routing rules

- **Route**: sets a default route to use wg0 and _table 2468_
- **RoutingPolicyRule**: unmarked packets use _table 2468_
- **RoutingPolicyRule**: bypass wireguard / _table 2468_ if using non default route

results:

**umarked packet**: -> _table 2468_ -> _wg0_ -> marked -> table main

```ini
[Match]
Name = wg0

[Network]
Address = CLIENT_IP_ADDRESS/SUBNET

[Route]
Destination = 0.0.0.0/0
Table = 2468

[RoutingPolicyRule]
InvertRule = true
FirewallMark = 1234
Table = 2468

[RoutingPolicyRule]
Table = main
SuppressPrefixLength = 0
```
