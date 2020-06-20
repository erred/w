---
description: using TURN relay as a proxy
title: TURN proxy
---

### _TURN_

Research [project](https://github.com/seankhliao/uva-rp2) for June.

NATs and firewalls are the bane of peer 2 peer (p2p) connections.

terms:

- _STUN_: Session Traversal Utilities for NAT, base protocol, used to learn your own external address
- _TURN_: Traversal Using Relay NAT, fallback when holepunching doesn't work, relay data through a server
- _ICE_: Interactive Connectivity Establishment, framework that offers a "complete" solution

#### TURN basics

base protocol, 3 different ways, control and data reuse the same socket pairs

```txt
Client              TURN Relay            Peer
  |                      |                  |
  |   UDP(TURN(data))    |     UDP(data)    |
  | -------------------> | ---------------> |
  |                      |                  |

Client              TURN Relay            Peer
  |                      |                  |
  |   TCP(TURN(data))    |     UDP(data)    |
  | -------------------> | ---------------> |
  |                      |                  |

Client              TURN Relay            Peer
  |                      |                  |
  | TCP(TLS(TURN(data))) |     UDP(data)    |
  | -------------------> | ---------------> |
  |                      |                  |
```

##### rfc6062, TCP to peer

```txt
Client              TURN Relay            Peer
  |                      |                  |
  |  TCP(TURN(control))  |                  |
  | -------------------> |                  |
  |   TCP(TURN(data))    |     TCP(data)    |
  | -------------------> | ---------------> |
  |                      |                  |

Client                 TURN Relay            Peer
  |                         |                  |
  | TCP(TLS(TURN(control))) |                  |
  | ----------------------> |                  |
  |  TCP(TLS(TURN(data)))   |     TCP(data)    |
  | ----------------------> | ---------------> |
  |                         |                  |
```

#### _Proxy_ forwarding

exposing a SOCKS5 interface

uses,

- get a generic socks proxy?
- pivot into the relay's private net

##### forwarding UDP

Proxy-Relay can also use TCP/TLS

```txt
Client                Proxy           TURN Relay         Peer
  |                     |                 |               |
  | TCP(SOCKS(control)) |                 |               |
  | ------------------> |                 |               |
  |   UDP(SOCKS(data))  | UDP(TURN(data)) |   UDP(data)   |
  | ------------------> | --------------> | ------------> |
  |                     |                 |               |
```

##### forwarding TCP

Proxy-Relay can also use TLS

```txt
Client                Proxy             TURN Relay          Peer
  |                     |                    |               |
  |                     | TCP(TURN(control)) |               |
  |                     | -----------------> |               |
  |   TCP(SOCKS(data))  |  TCP(TURN(data))   |   TCP(data)   |
  | ------------------> | -----------------> | ------------> |
  |                     |                    |               |
```

#### reverse shell

Proxy Reverse - Relay can also use TCP/TLS

udp

```txt
Target     Proxy Reverse             TURN Relay       Proxy Server               Client
  |              |                        |                  |                     |
  |              | UDP(TURN(QUIC(hello))) | UDP(QUIC(hello)) |                     |
  |              | ---------------------> |----------------> |                     |
  |              |                        |                  | TCP(SOCKS(control)) |
  |              |                        |                  | <------------------ |
  |   UDP(data)  | UDP(TURN(QUIC(data)))  | UDP(QUIC(data))  |  UDP(SOCKS(data))   |
  | <----------- | <--------------------- | <--------------- | <------------------ |
```

tcp

```txt
Target     Proxy Reverse             TURN Relay       Proxy Server            Client
  |              |                        |                  |                  |
  |              | UDP(TURN(QUIC(hello))) | UDP(QUIC(hello)) |                  |
  |              | ---------------------> |----------------> |                  |
  |   TCP(data)  | UDP(TURN(QUIC(data)))  | UDP(QUIC(data))  | TCP(SOCKS(data)) |
  | <----------- | <--------------------- | <--------------- | <--------------- |
```

#### _Problems_

- TCP to peer support is rare, in both relays (only Coturn) and libraries (self implemented)
- TURN credentials can be hard to extract, only long term credentials can allocate, short term creds for ICE can't
- Wireshark and browser WebRTC don't mix
