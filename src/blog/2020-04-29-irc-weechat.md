---
description: irc with weechat
title: irc weechat
---

### _irc_ weechat

some developer communities apparently only exist on irc,
so you need a proper client to connect.

todo: make weechat less ugly

#### _concepts_

- _server:_ thing you connect to, doesn't keep history
- _username:_ useless field, everything uses nick(s)
- _password:_ authentication for login
- _nick:_ alternate names you can use
- _nickserv:_ registration of nicks you can use
- _sasl:_ login/auth protocol

since servers don't keep global history,
you need to be connected to be able to see history.

- _bouncer:_ middleware client so you are always connected to the server
- _weechat:_ a terminal client, can also be a bouncer

#### setup goals

- registered nick on freenode
- weechat on a server as bouncer
- weechat on laptop as client
- weechat app on android phone as client

##### server

connect to freenode

```
# /server add name-of-server server-domain/port -options
/server add freenode chat.freenode.net/6697 -ssl -autoconnect
/connect freenode
```

register a name

```
# change to name you want to use
/nick name-you-want
/msg NickServ REGISTER your-password your@email.com
# follow instructions in email
```

setup login

```
/secure set name-of-secret your-password
/set irc.server.freenode.sasl_username name-you-want
/set irc.server.freenode.sasl_password ${sec.data.name-of-secret}
```

```
/join #channel-you-want-to-join
/set irc.server.freenode.autojoin #channels,#you,#always,#want,#to,#join
```

setup relay (allow clients to connect)

expects `cat privkey.pem fullchain.pem` combined cert in `$WEECHAT_HOME/ssl/relay.pem`

```
/secure set another-secret relay-password
/relay add ssl.irc port-for-irc-over-ssl
/relay add ssl.weechat port-for-weechat-over-ssl
/set relay.network.password ${sec.data.another-secret}
```

save and start in a screen / tmux / headless session in systemd

```
/save
```

#### laptop client

```
# /server add name-of-relay host-of-relay/port-for-irc-over-ssl -ssl
/server add weechat weechat.example.com/7992 -ssl -autoconnect

# /secure set yet-another-secret server:relay-password
/secure set yet-another-secret freenode:relay-password
/set irc.server.weechat.password ${sec.data.yet-another-secret}
/set irc.server.weechat.nicks name-you-want
```

#### mobile client

top right > settings > connection

```
Connection type: WeeChat SSL
Relay host: host-of-relay
Relay port: port-for-weechat-over-ssl
Relay password: relay-password
```
