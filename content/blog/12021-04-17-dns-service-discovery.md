---
title dns service discovery
description: who reads RFCs...
---

### _dns_ service discovery

Globally available key-value store,
let's use it for service discovery!

iana service name and port
[directory](https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.txt)

#### _SRV_

Basic: a single instance per service.
(can have multiple records, but each host is expected to serve the same data)

```txt
_service._protocol.zone. SRV "priority weight port instance1.zone"
```

#### _DNS-SD_

Turns the `_service._protocol` into a listing of instances.

```txt
_service._protocol.zone. PTR instance1._service.protocol.zone.
_service._protocol.zone. PTR instance2._service.protocol.zone.
```

and each instance gets an extra `TXT` to store additional data needed per protocol.

```txt
instance1._service.protocol.zone. SRV "priority weight port instance1.zone"
instance1._service.protocol.zone. TXT "key1=value1,key2=value2"

instance2._service.protocol.zone. SRV "priority weight port instance2.zone"
instance2._service.protocol.zone. TXT "key1=value1,key2=value2"
```

#### _how_ can i use this

Totally not tested. On DNS server

```txt
_wireguard._udp.example.com. PTR base32pubkey1._wireguard._udp.example.com.

base32pubkey1._wireguard._udp.example.com. SRV "10 0 41630 mylaptop.local"
base32pubkey1._wireguard._udp.example.com. SRV "20 0 ????? mylaptop.wgsd.example.com."
base32pubkey1._wireguard._udp.example.com. TXT ""
```

and have `mylaptop` respond to mDNS for local lookup,
falling back to something like [wgsd](https://github.com/jwhited/wgsd).
