### _weechat_ config updates

use it more, config it more

#### bouncer mode

configure at least the important parts:
server autojoin, secured data, relay

generate a self signed ssh certificate,
adjust `DNS:...,IP:...` as appropriate

uses ecdsa key,
ed25519 certs are apparently a bit too exotic

```sh
openssl req -x509 -nodes -days 7300 \
  -subj "/O=weechat/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:0.0.0.0" \
  -newkey ec:<(openssl ecparam -name prime256v1) \
  -keyout relay.pem -out relay.pem

# ed25519?
openssl req -x509 -nodes -days 7300 \
  -key <(openssl genpkey -algorithm ED25519) \
  -subj "/O=weechat/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:0.0.0.0" \
  -keyout relay.pem -out relay.pem
```

obtain fingerprint with

```sh
openssl x509 -in relay.pem -outform der | sha256sum -b | cut -d ' ' -f1
```

add systemd user service (from archwiki)

```systemd
[Unit]
Description=A headless WeeChat client and relay service
After=network.target

[Service]
Type=forking
ExecStart=/usr/bin/weechat-headless --daemon

[Install]
WantedBy=default.target
```

enable service and enable linger for user
(run service on boot, don't stop when logout),
reboot

```sh
systemctl --user daemon-reload
systemctl --user enable weechat-headless
sudo loginctl enable-linger $USER
sudo reboot
```

#### weechat configs

in alphabetical order,
but secured data should probably be set first

##### buflist

hide all the merged server buffer

```txt
/set buflist.look.display_conditions "${buffer.hidden}==0 && ${merged}==0 || ${name}==weechat"
```

##### irc

stop jumping buffer

```txt
/set irc.look.buffer_switch_autojoin off
```

more colors

```txt
/set irc.look.color_nicks_in_names on
/set irc.look.color_nicks_in_nicklist on
```

disable ctcp

```txt
/set irc.ctcp.action ""
/set irc.ctcp.clientinfo ""
/set irc.ctcp.finger ""
/set irc.ctcp.ping ""
/set irc.ctcp.source ""
/set irc.ctcp.time ""
/set irc.ctcp.userinfo ""
/set irc.ctcp.version ""
```

default messages

```txt
/set irc.server_default.msg_part "going... going... gone!"
/set irc.server_default.msg_quit "it can't be DNS?! it's always DNS!!!"
/set irc.server_default.nicks ""
/set irc.server_default.username ""
```

default username

```txt
/set irc.server_default.username seankhliao
```

server for bouncer

```txt
/server add bouncer 192.168.100.1/7992
/set irc.server.bouncer.ssl on
/set irc.server.bouncer.ssl_fingerprint = "6064898BD86791CF681593BD41B86541F2EFE6B34708D95C1ED7412794626528"
/set irc.server.bouncer.password = "freenode:${sec.data.relay_pass}"
/set irc.server.bouncer.autoconnect on
```

server for freenode

```txt
freenode.addresses = "chat.freenode.net/6697"
freenode.ssl = on
freenode.sasl_mechanism = plain
freenode.sasl_username = "seankhliao"
freenode.sasl_password = "${sec.data.freenode_pass}"
freenode.autoconnect = on
freenode.autojoin = "..."
```

###### CertFP

login and identify with client certs

generate cert as above with name `freenode.pem`, add

```txt
/set irc.server.freenode.ssl_cert %h/ssl/freenode.pem
/set irc.server.freenode.sasl_mechanism external
```

check and associate current cert with identity,
future logins won't require password

```txt
/whois YourOwnNick
/msg NickServ CERT ADD
```

##### logger

don't log join/leave, and dim log line

```txt
/set logger.color.backlog_end darkgray
/set logger.color.backlog_line darkgray
/set logger.file.color_lines on
/set logger.level.irc 2
```

##### scripts

colorize_nicks

```txt
/script install colorize_nicks
/set colorize_nicks.look.colorize_input on
```

##### relay

open up a protected relay

```txt
/set relay.network.password ${sec.data.relay_pass}
/set relay.irc.backlog_since_last_message on
/set relay.port.ssl.irc 7992
/set relay.port.ssl.weechat 7993
```

unlimited backlog

```txt
/set relay.irc.backlog_max_number 0
```

##### secure

store data "securely",
decrypt automatically with passphrase in a file

```txt
/secure passphrase
/set sec.crypt.passphrase_file ~/.ssh/weechat_pass
/secure set relay_pass xxxxxx
/secure set freenode_pass yyyyyy
```

##### weechat

disable the giant logo

```txt
/set weechat.startup.display_logo off
```

change the look of the left gutter:
more compact time,
change icons

```txt
/set weechat.look.buffer_time_format  "${color:253}%H${color:245}%M"
/set weechat.look.prefix_action "*"
/set weechat.look.prefix_error "!!"
/set weechat.look.prefix_join ">>"
/set weechat.look.prefix_network "~"
/set weechat.look.prefix_quit "<<"
/set weechat.look.read_marker_string "*"
/set weechat.look.read_marker_always_show on
```

colors:
_purple_ (5) is me,
more colors for everyone else

```txt
/set weechat.color.chat_delimiters = darkgray
/set weechat.color.chat_highlight = 5
/set weechat.color.chat_highlight_bg = 0
/set weechat.color.chat_nick_colors = "1,2,3,4,6,7,9,10,11,12,13,14,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,69,70,182,183,184,224,225,226,227"
/set weechat.color.chat_nick_self = 5
/set weechat.color.chat_prefix_suffix = gray
```

disable unused plugins

```txt
/set weechat.plugin.autoload = "*,!fifo,!guile,!javascript,!lua,!perl,!php,!ruby,!tcl,!spell,!xfer"
```

move buffers to top,
who looks at the nicklist with 1800 people in a channel,
abbreviated statusline,
no titlebar

```txt
/set weechat.bar.buflist.position = top
/set weechat.bar.nicklist.hidden = on
/set weechat.bar.status.color_bg = 0
/set weechat.bar.status.items = "+buffer_name+(buffer_modes)+{buffer_nicklist_count}+buffer_zoom+buffer_filter,scroll,[lag],[hotlist],completion"
/set weechat.bar.title.hidden = on
```

filter out join/leaves, except for recent partiticants,
equivalent to
`/filter add irc_smart * irc_smart_filter *`

```txt
/set irc.filter.irc_smart on
/filter add irc_smart * irc_smart_filter *
```

navigate with ctrl+arrows

```txt
/set weechat.key.meta2-1;5A = "/window page_up"
/set weechat.key.meta2-1;5B = "/window page_down"
/set weechat.key.meta2-1;5C = "/buffer +1"
/set weechat.key.meta2-1;5D = "/buffer -1"
```
