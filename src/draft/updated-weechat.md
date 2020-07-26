### _weechat_ config updates

use it more, config it more

#### weechat

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
```

colors:
purple is me,
more colors for everyone else

```txt
/set weechat.color.chat_delimiters = darkgray
/set weechat.color.chat_highlight = 5
/set weechat.color.chat_highlight_bg = 0
/set weechat.color.chat_nick_colors = "1,2,3,4,6,7,9,10,11,12,13,14,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,69,70,182,183,184,224,225,226,227"
/set weechat.color.chat_nick_self = 5
/set weechat.color.chat_prefix_suffix = gray
```

more history

```txt
/set weechat.history.max_buffer_lines_number = 1024
```

disable the plugins I won't use

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
/set irc.filter.irc_smart on;*;irc_smart_filter;*
```

navigate with ctrl+arrows

```txt
/set weechat.key.meta2-1;5A = "/window page_up"
/set weechat.key.meta2-1;5B = "/window page_down"
/set weechat.key.meta2-1;5C = "/buffer +1"
/set weechat.key.meta2-1;5D = "/buffer -1"
```

#### buflist

list of open buffers, hide all the merged server ones

```txt
/set buflist.look.display_conditions "${buffer.hidden}==0 && ${merged}==0 || ${name}==weechat"
```

#### logger

don't log join/leave, and dim log line

```txt
/set logger.color.backlog_end darkgray
/set logger.color.backlog_line darkgray
/set logger.file.color_lines on
/set logger.level.irc 2
```

#### secure

#### relay

open up a protected relay

```txt
/set relay.network.password ${sec.data.relay}
/set relay.irc.backlog_since_last_message on
/set relay.port.ssl.irc 7992
/set relay.port.ssl.weechat 7993
```

#### irc

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
