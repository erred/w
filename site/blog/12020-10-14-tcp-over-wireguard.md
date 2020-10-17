---
description: running tcp over wireguard
title: tcp over wireguard
---

something I noticed:
less disconnects (but more timeouts) when running things over over wireguard.

possible reason?
Most of the network instability is in the last hop,
wifi router - laptop or cell tower - phone.
This means the device connects/disconnects will affect the availability of an ip address
and by extension, an associated connections.
This is masked by the always available nature of wireguard,
which prefers to send your packets into the void.
