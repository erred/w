---
description: layers and layers of protocols
title: multiplexing protocols
---

### _multiplexing_

_hardware_ layer: you have 2 interfaces:
no need for addressing, just dump everything on the wire

_link_ layer: you have more than 2 interfaces all connected together:
give each interface an address (MAC address, random assignment)

_internet_ layer: you have a lot more than 2 interfaces all connected together:
give each interface an address (IP address, heirarchical assignment)

_transport_ layer: you want to talk to a specific process listening on an interface:
give each process a port, optionally handle streaming (ex: UDP/TCP port)

_transport_ layer II: you want to talk to a specific resource set in a process (and/or you want security):
give each resource set a domain name (ex: DTLS/TLS/QUIC authority, SNI)

_application_ layer: you want to talk to a specific resource in your resource set:
give each resource a path, optionally duplicate TLS authority, optionally add method (ex: HTTP path)

_application_ layer II: whatever you built on top of HTTP:
(ex: JSON-RPC, gRPC, GrapQL, ...)

#### _so..._

ex raw DNS works directly on top of the transport layer, has no concept of host/domain names,
`drill @your.resolver query` does a lookup first for you resolver where the query is then sent,
so you can't really multiplex multiple DNS resolvers on a single port.
