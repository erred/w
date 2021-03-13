---
description: modernish takes on the cli tools
title: modern cli toolbox
---

### _cli_ tools

There are 1800 things in my `/bin`,
I'm sure some of them can do with an update.

#### _replacements_

Modern takes, usually with colors or expanded functionality

| new            | old              | description                                |
| -------------- | ---------------- | ------------------------------------------ |
| [bat][bat]     | cat, less        | colored pager                              |
| [exa][exa]     | ls, tree         | ls with git features                       |
| [mtr][mtr]     | traceroute, ping | pager updating ping times per hop          |
| [ncdu][ncdu]   | du, dust         | curses disk usage viewer                   |
| [neovim][nvim] | vi, vim          | a better vim? going all in on lua          |
| [ripgrep][rg]  | grep             | faster grep                                |
| [rsync][rsync] | scp              | differential file sync                     |
| [sd][sd]       | sed              | more intuitive sed                         |
| [skim][sk]     | locate, fzf      | interactive search, filenames and contents |
| [cw][cw]       | wc               | faster wc                                  |
| [fd][fd]       | find             | more intuitive find                        |

#### _newish_

maybe not the first to implement an idea, but good ideas/implementation nonetheless

| tool             | description                                         |
| ---------------- | --------------------------------------------------- |
| [entr][entr]     | run command on file change                          |
| [tag][tag]       | generate terminal shortcuts to ripgrep results      |
| [tokei][tokei]   | count lines of code                                 |
| [xsv][xsv]       | CSV toolkit                                         |
| [jq][jq]         | JSON toolkit                                        |
| [rclone][rclone] | unified cli to cloud storage                        |
| [aria2][aria2]   | multiprotocol downloader (http, ftp,bittorrent,...) |

[aria2]: https://github.com/aria2/aria2
[bat]: https://github.com/sharkdp/bat
[entr]: https://github.com/eradman/entr/
[exa]: https://github.com/ogham/exa
[mtr]: https://github.com/traviscross/mtr
[ncdu]: https://dev.yorhel.nl/ncdu
[nvim]: https://github.com/neovim/neovim
[rg]: https://github.com/BurntSushi/ripgrep
[rsync]: https://rsync.samba.org/
[sd]: https://github.com/chmln/sd
[sk]: https://github.com/lotabout/skim
[tag]: https://github.com/aykamko/tag
[tokei]: https://github.com/XAMPPRocky/tokei
[xsv]: https://github.com/BurntSushi/xsv
[jq]: https://github.com/stedolan/jq
[rclone]: https://github.com/rclone/rclone
[cw]: https://github.com/Freaky/cw
[fd]: https://github.com/sharkdp/fd
