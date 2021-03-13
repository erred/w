---
description: finally getting round to remembering iptables
title: iptables notes
---

### iptables basics

_FILTER_ is the default table,
traffic passes through 1 of 3 chains { _INPUT_, _FORWARD_, _OUTPUT_ }
(most of the time)

container traffic passes through _FORWARD_

_PREROUTING_, _POSTROUTING_ only available on _RAW_, _MANGLE_

```txt
PREROUTING ─┬──── FORWARD ────┬─ POSTROUTING
            │                 │
      INPUT │                 │ OUTPUT
            │                 │
            └─ local process ─┘
```

#### filter traffic

chain modifiers:

- _-P_: set chain default policy
- _-F_: flush all rules
- _-A_: append rule
- _-D_: delete rule

selectors:

- _-i_: interface
- _-s_ / _-d_: source/destination address range
- _-p_: protocol
- _--sport_ / _--dport_: source/destination port (for tcp/udp)

actions:

- _-j_: jump to target
- _ACCEPT_: allow packet to pass
- _DROP_: drop packet
- _RETURN_: continue processing from last jump
- _REJECT_: reject packet, ex: host/network unreachable

```sh
iptables -P FORWARD ACCEPT
iptables -A FORWARD -s 0.0.0.0/0 -d 1.2.3.4/32 -p tcp --dport 80 -j ACCEPT
```

#### allow established connections through

- _-m_: match using an extention module

```sh
iptables -A INPUT -m conntrack --cstate RELATED,ESTABLISHED -j ACCEPT
iptables -A OUTPUT -m conntrack --cstate RELATED,ESTABLISHED -j ACCEPT
iptables -A FORWARD -m conntrack --cstate RELATED,ESTABLISHED -j ACCEPT
```
