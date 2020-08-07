---
description: updated weechat configs
title: updated weechat
---

### _weechat_ config updates

use it more, config it more

#### _bouncer_ mode

configure at least the important parts:
server login/autojoin, secured data, relay

generate a self signed ssh certificate,
adjust `DNS:...,IP:...` as appropriate

```sh
openssl req -x509 -nodes -days 7300 \
  -key <(openssl genpkey -algorithm ED25519 | tee relay.pem) \
  -subj "/O=weechat/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:0.0.0.0" \
  -keyout relay.pem -out relay.pem

# or as 2 separate calls if you have weird issues
openssl genpkey -algorithm ED25519 -out relay.pem
openssl req -x509 -nodes -days 7300 \
  -key relay.pem \
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

#### weechat _configs_

in alphabetical order,
but secured data should probably be set first

##### _buflist_

hide all the merged server buffer

```txt
/set buflist.look.display_conditions ${buffer.hidden}==0 && ${merged}==0 || ${name}==weechat
```

##### _irc_

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
/set irc.server_default.msg_part going... going... gone!
/set irc.server_default.msg_quit it can't be DNS?! it's always DNS!!!
/set irc.server_default.nicks ""
/set irc.server_default.username ""
```

###### _bouncer_ server

```txt
/server add bouncer 192.168.100.1/7992
/set irc.server.bouncer.ssl on
/set irc.server.bouncer.ssl_fingerprint 33300c58504afc85a43ab49acea568d37780b214b078f8917eb9033a762b8e26
/set irc.server.bouncer.autoconnect on
/set irc.server.bouncer.username seankhliao
/set irc.server.bouncer.nicks seankhliao,arccy
/set irc.server.bouncer.password freenode:${sec.data.relay_pass}
```

###### _freenode_ server

```txt
/server add freenode chat.freenode.net/6697
/set irc.server.freenode.ssl on
/set irc.server.freenode.autoconnect on
/set irc.server.freenode.autojoin #archlinux,#archlinux-offtopic,#archlinux-reproducible,#go-nuts
/set irc.server.freenode.username seankhliao
/set irc.server.freenode.nicks seankhliao,arccy
```

###### _SASL_ with password

login and identify with password

```txt
/set irc.server.freenode.sasl_mechanism plain
/set irc.server.freenode.sasl_username seankhliao
/set irc.server.freenode.sasl_password ${sec.data.freenode_pass}
```

###### _SASL_ with CertFP

login and identify with client certs

generate cert as above with name `freenode.pem`, add

```txt
/set irc.server.freenode.ssl_cert %h/ssl/freenode.pem
/set irc.server.freenode.sasl_mechanism external
```

check and associate current cert with identity,
future logins won't require password

```txt
/msg NickServ identify PASSWORD
/whois seankhliao
/msg NickServ CERT ADD
```

##### _logger_

don't log join/leave, and dim log line

```txt
/set logger.color.backlog_end darkgray
/set logger.color.backlog_line darkgray
/set logger.file.color_lines on
/set logger.level.irc 2
```

##### _scripts_

colorize_nicks

```txt
/script install colorize_nicks.py
/set colorize_nicks.look.colorize_input on
/set colorize_nicks.look.min_nick_length 3
```

##### _relay_

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

##### _secure_

store data "securely",
decrypt automatically with passphrase in a file

```txt
/secure passphrase
/set sec.crypt.passphrase_file ~/.ssh/weechat_pass
/secure set relay_pass xxxxxx
/secure set freenode_pass yyyyyy
```

##### _trigger_

stylized messages

```txt
/trigger del beep
/trigger addreplace irc_join modifier "2000|weechat_print" "${tg_tags} =~ ,irc_join," "/.*[^(]\((.*)\).*/${color:237}${tg_tag_nick}@${re:1}/tg_message_nocolor /.*/${tg_prefix}\t${tg_message_nocolor}"
/trigger addreplace irc_nick modifier "2000|weechat_print" "${tg_tags} =~ ,irc_nick," "/.*irc_nick1_([^,]*),irc_nick2_([^,]*).*/${re:1} -> ${re:2}/tg_tags /.*/${tg_prefix}\t${tg_tags}"
/trigger addreplace irc_quit modifier "2000|weechat_print" "${tg_tags} =~ ,irc_quit," "/.*[^(]\((.*)\).*\ (\(.*\))/${color:237}${tg_tag_nick}@${re:1}: ${re:2}/tg_message_nocolor /.*/${tg_prefix}\t${tg_message_nocolor}"
/trigger addreplace irc_part modifier "2000|weechat_print" "${tg_tags} =~ ,irc_part," "/.*[^(]\((.*)\).*\ (.*)/${color:237}${tg_tag_nick}@${re:1}: ${re:2}/tg_message_nocolor /.*/${tg_prefix}\t${tg_message_nocolor}"
/trigger addreplace nick_color_action modifier weechat_print "${tg_tags} =~ ,irc_action, && ${tg_tags} !~ ,self_msg," "/.*/${info:nick_color,${tg_tag_nick}}${tg_prefix_nocolor}\t${tg_message}"
/trigger addreplace url_color modifier weechat_print "${tg_tags} !~ irc_quit" ";[a-z]+://\S+;${color:32}${re:0}${color:reset};"
```

from [weechat wiki](https://github.com/weechat/weechat/wiki/Triggers),
dim user

```txt
/alias add dim trigger addreplace dim_$server_$1 modifier weechat_print "${tg_tag_nick} == $1 && \${server} == $server" "/(.*)/${color:darkgray}${tg_prefix_nocolor}\t${color:darkgray}${tg_message_nocolor}/"; print \---\t$1 is now dimmed on $server
/alias add undim trigger del dim_$server_$1; print \---\t$1 is no longer dimmed
```

##### _weechat_

disable the giant logo

```txt
/set weechat.startup.display_logo off
```

change the look of the left gutter:
more compact time,
change icons

```txt
/set weechat.look.buffer_time_format  ${color:253}%H${color:245}%M
/set weechat.look.prefix_action *
/set weechat.look.prefix_error !!
/set weechat.look.prefix_join >>
/set weechat.look.prefix_network ~
/set weechat.look.prefix_quit <<
/set weechat.look.read_marker_string *
```

colors:
_purple_ (5) is me,
more colors for everyone else

```txt
/set weechat.color.chat_delimiters darkgray
/set weechat.color.chat_highlight 5
/set weechat.color.chat_highlight_bg 0
/set weechat.color.chat_nick_colors 1,2,3,4,6,7,9,10,11,12,13,14,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,69,70,182,183,184,224,225,226,227
/set weechat.color.chat_nick_self 5
/set weechat.color.chat_prefix_suffix gray
```

disable unused plugins

```txt
/set weechat.plugin.autoload *,!fifo,!guile,!javascript,!lua,!perl,!php,!ruby,!tcl,!spell,!xfer
```

move buffers to top,
who looks at the nicklist with 1800 people in a channel,
abbreviated statusline,
no titlebar

```txt
/set weechat.bar.buflist.position top
/set weechat.bar.nicklist.hidden on
/set weechat.bar.status.color_bg 0
/set weechat.bar.status.items +buffer_name+(buffer_modes)+{buffer_nicklist_count}+buffer_zoom+buffer_filter,scroll,[lag],[hotlist],completion
/set weechat.bar.title.hidden on
```

filter out join/leaves, except for recent partiticants,

```txt
/set irc.look.smart_filter on
/filter add irc_smart * irc_smart_filter *
```

navigate with ctrl+arrows

```txt
/key bind meta2-1;5A /window page_up
/key bind meta2-1;5B /window page_down
/key bind meta2-1;5C /buffer +1
/key bind meta2-1;5D /buffer -1
```
