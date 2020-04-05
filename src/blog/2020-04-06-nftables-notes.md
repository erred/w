---
description: so iptables is getting replaced
title: nftables notes
---

### nftables basics

no default anythings

_tables_ (per address family) hold _chains_ hold _rules_

```sh
nft add table inet table_name
nft list tables

nft add chain inet table_name chain_name
nft list table inet table_name
```

_chains_ are triggered by _hooks_ or other chains,
all 5 chains available for _ip_ / _ip5_ / _inet_ / _bridge_

_chains_ have types (with self explanatory names): _filter_, _nat_, _route_

```txt
PREROUTING ─┬──── FORWARD ────┬─ POSTROUTING
            │                 │
      INPUT │                 │ OUTPUT
            │                 │
            └─ local process ─┘
```

```sh
nft add chain inet table_name chain_name '{ type filter hook input priority 0 policy accept; }'
```

_rules_

- saddr/daddr
- tcp/udp
- sport/dport

```sh
nft add rule inet table_name chain_name saddr 0.0.0.0/0 daddr 1.2.3.4/32 tcp dport 80 accet
```

_sets_

named, reuseable, updateable set of addresses,
use with `@set_name`

```sh
nft add set inet table_name set_name '{ type ipv4_addr }'
nft add element inet table_name set_name '{ 1.2.3.4 }'
```

### allow established connections

```sh
nft add rule inet table_name chain_name ct state related,established accept
```

### config file

```txt
table inet table_name {
  chain chain_name {
    type filter hook input priority 0;

    ct state {established,related} accept
    ct state invalid drop

    tcp dport 80 accept
  }
}
```
