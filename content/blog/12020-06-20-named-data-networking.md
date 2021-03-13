---
description: named data networking overview
title: named data networking
---

### _NDN_ Named Data Networking

Research [project](https://github.com/erred/uva-rp1) topic from January.
Finally getting round to writing it down.

#### _Do_ I need to know about it

NDN is an academic curiosity
and is unlikely to progress beyond that.

_IPFS_ is better.

#### _Current_ Internet

- _IP_: addresses / routes between hosts (machines)
- _TCP_ / _UDP_ / _..._: addresses processes (applications)
- _HTTP_ / _FTP_ / _..._: addresses data (content)

#### _Envisioned_ NDN Internet

Content is addressed directly,
hosts / processes either have and serve it,
or requests it from upstreams and serves it when it in turn.
Content has heirarchical names, and routed according to name heirarchy,
basically giant global distributed filesystem.

Names (and caching) is permanent,
but not enforced algorithmically (ex. hash).
Everything is Type-Length-Value (TLV) encoded (but inconsistently),.
Goal would be to replace the IP/TCP/HTTP (or equivalent) stack,
implementing routers on hardware.

#### _Reality_ of NDN

Unsurprisingly nobody wants to invest in implementing hardware
for an inefficient protocol.
The only current [implementation](https://github.com/named-data/NFD)
works primarily on top of _UDP_, optionally Ethernet, TCP, WebSockets.
HTTP can be implemented on top of NDN.

So you can have a Ethernet / TCP / HTTP / WebSocket / NDN / HTTP stack.
