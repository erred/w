---
title: jool nat64, coredns dns64
description: setting up nat64 and dns64
---

### _ipv6_

Why can't the world just use ipv6?
Ah yes, longer addresses that are hard to remember,
another set of dns entries to create,
and sometimes (often?) your users don't have it anyway so why bother?

Anyway, you have some ipv6 network (maybe k8s kind in ipv6 mode?)
that needs to reach out to the ipv4 internet.
You (usually) need 2 things: _nat64_ and _dns64_.

#### _dns64_

Problem: you request `example.com` but it only has an `A` record.

Solution: reply with a synthesized `AAAA` record using mapped addresses
(a section of ipv6 address space mapping directly to ipv4, eg `64:ff9b::192.0.2.128`).
[RFC 6052](https://datatracker.ietf.org/doc/html/rfc6052) even reserves
the `64:ff9b::/96` range for translation.

##### _coredns_

Most documentation online uses [bind](https://www.isc.org/bind/),
but whatever, we're using [coredns](https://coredns.io/).
Enable the `dns64` plugin, and `forward` everything to upstream.
Oh, and set `acl` so you don't end up running an open resolver
(if you run it somewhere with public addresses).
This doesn't have to run on a dual stack node,
it just needs to talk to an upstream.

```Corefile
. {
	dns64
	forward . 8.8.8.8:53

	errors
	log
	acl {
		allow net 10.0.0.0/8 172.16.0.0/12 192.168.0.0/16 127.0.0.0/8
		allow net fc00:f853:ccd:e793::/64 ::1/128
		block
	}
}
```

#### _nat64_

Problem: you have your "fake" ipv6 address for your destination now,
but you still need to be able to send data to it.

Solution: have your gateway (dual stack) understand your mapped prefix
and translate + nat everything with a destination there from ipv6 to ipv4.

##### _jool_

[jool](https://www.jool.mx/en/index.html) appears to be the only reasonable choice here.
It's a kernel module. Install it, `modprobe jool`, `jool file handle $config`.
`pool6` is the prefix to recognize as containing ipv4 addresses.
One thing to note, it's [implemented as `PRE_ROUTING`](https://www.jool.mx/en/faq.html#why-is-my-ping-not-working)
so you need to have a separate network namespace to test it.

```json
{
	"instance": "default",
	"framework": "netfilter",
	"global": {
		"pool6": "64:ff9b::/96"
	}
}
```
