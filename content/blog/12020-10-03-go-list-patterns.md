---
description: go list pattern coverage
title: go list patterns
style: |
  picture {
    margin-top: 0;
  }
---

### _list_

The "I know it lists things but I don't remember what" command.

reference repo: [erred/go-list-ex](https://github.com/erred/go-list-ex)

#### _dependency_ graph

This is everything:

![graph of everything](https://raw.githubusercontent.com/erred/go-list-ex/main/base.svg)

#### _module_ mode

listing modules

##### go list -m

![go list -m](https://raw.githubusercontent.com/erred/go-list-ex/main/m.svg)

##### go list -m all

![go list -m all](https://raw.githubusercontent.com/erred/go-list-ex/main/mall.svg)

#### _package_ mode

##### go list ./...

![go list ./...](https://raw.githubusercontent.com/erred/go-list-ex/main/dotdotdot.svg)

##### go list -deps ./...

![go list -deps ./...](https://raw.githubusercontent.com/erred/go-list-ex/main/dotdotdotdeps.svg)

##### go list -deps -test ./...

![go list -deps -test ./...](https://raw.githubusercontent.com/erred/go-list-ex/main/dotdotdotdepstest.svg)

##### go list all

_all_ changed meaning between 1.15 and 1.16

###### go1.16+

![1.16 go list all](https://raw.githubusercontent.com/erred/go-list-ex/main/all.svg)

###### go1.15 or earlier

![1.15 go list all](https://raw.githubusercontent.com/erred/go-list-ex/main/all115.svg)
