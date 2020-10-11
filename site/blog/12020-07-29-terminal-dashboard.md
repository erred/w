---
description: a all in one dashboard for the terminal
title: terminal dashboard
---

### _terminal_ dashboard

sometimes you just want to log in, hit a dashboard and see wtf your server is doing

#### _options_

_tldr_: ytop or bashtop for all in one, htop for process tree

- [atop][atop]
  - \+ cpu, memory, io, network, process, disk
  - \- no process tree
  - ~ C?
  - weird ui
- [bashtop][bashtop]
  - \+ cpu, memory, network, process, disk, process tree
  - \- io
  - ~ bash
  - looks sleak, slow startup
- [bmon][bmon]
  - only does network
- [bottom][bottom]
  - \+ cpu, memory, io, network, process, disk
  - \- no process tree
  - \+ rust
  - impossible to find
- [glances][glances]
  - \+ cpu, memory, io, network, process, disk
  - \- no process tree
  - \- python
- [gotop][gotop]
  - \+ cpu, memory, io, network, process, disk
  - \- no process tree
  - \+ go
- [gtop][gtop]
  - \- javascript
- [htop][htop]
  - \+ cpu, memory, io, proces, process tree
  - \- network, disk
  - \~ c
- [iftop][iftop]
  - only does network connections
- [iotop][iotop]
  - only does io
- [vtop][vtop]
  - \+ cpu, memory, proces
  - \- io, network, disk
  - \- no process tree
  - \- javascript
- [wtf][wtf]
  - customizable dashboard, doesn't do monitoring out of the box
- [ytop][ytop]
  - \+ cpu, memory, io, network, process, disk
  - \- no process tree
  - \+ rust
  - same author as gotop, exact same ui
- [zenith][zenith]
  - \+ cpu, memory, io, network, process, disk
  - \- no process tree
  - \+ rust
  - funky graphs

[atop]: https://www.atoptool.nl/
[bashtop]: https://github.com/aristocratos/bashtop
[bmon]: https://github.com/tgraf/bmon
[bottom]: https://github.com/ClementTsang/bottom
[glances]: https://github.com/nicolargo/glances
[gotop]: https://github.com/xxxserxxx/gotop
[gtop]: https://github.com/aksakalli/gtop
[htop]: https://github.com/hishamhm/htop
[iftop]: http://www.ex-parrot.com/pdw/iftop/
[iotop]: https://github.com/analogue/iotop
[vtop]: https://github.com/MrRio/vtop
[wtf]: https://github.com/wtfutil/wtf
[ytop]: https://github.com/cjbassi/ytop
[zenith]: https://github.com/bvaisvil/zenith
