---
title: bare git private go modules
description: someone complained it was too hard
---

### _go_ modules

`GOPATH` is dying, better switch to modules soon,
so how can you make it work with just a git (no http(s))?

#### _ssh_ setup

make `ssh` work automatically, ie `ssh arch.seankhliao.com`. `~/.ssh/config`:

```txt
Host arch.seankhliao.com
    User arccy                      # assuming you don't use a generic git user
    Port 443                        # assuming you use nonstandard port
    IdentityFile ~/.ssh/id_ecdsa_sk # ssh keyzzzz
```

#### _git_ setup

make `git` always use `ssh`, and point it to your home directory. `~/.gitconfig`:

```txt
[url "arccy@arch.seankhliao.com:"]
    insteadOf = "git://arch.seankhliao.com/"
    insteadOf = "https://arch.seankhliao.com/"
```

#### _dependency_

create a dependency!

```sh
$ mkdir -p ~/tmp/hello && cd ~/tmp/hello

$ cat << EOF > hello.go
package hello

var Hello = "world"
EOF
```

create a git repo on the remote

```sh
$ ssh arch.seankhliao.com git init --bare hello -b main
Initialized empty Git repository in /home/arccy/hello/
```

create a git repo locally and push to remote

```sh
$ git init -b main
Initialized empty Git repository in /home/arccy/tmp/hello/.git/

$ git add .

$ git commit -m "init"
[main (root-commit) 02adee5] init
 1 file changed, 3 insertions(+)
 create mode 100644 hello.go

$ git remote add origin https://arch.seankhliao.com/hello

$ git push -u origin main
Enumerating objects: 3, done.
Counting objects: 100% (3/3), done.
Writing objects: 100% (3/3), 240 bytes | 240.00 KiB/s, done.
Total 3 (delta 0), reused 0 (delta 0), pack-reused 0
To arch.seankhliao.com:hello
 * [new branch]      main -> main
Branch 'main' set up to track remote branch 'main' from 'origin'.
```

#### _main_

Now to use the dependency

```sh
$ mkdir -p ~/tmp/world && cd ~/tmp/world

# you may want to persist this setting with go env -w GOPRIVATE=...
$ export GOPRIVATE=arch.seankhliao.com

$ go mod init example/world
go: creating new go.mod: module example/world

$ cat << EOF > main.go
package main

import (
        "fmt"

        "arch.seankhliao.com/hello.git" // use .git to tell Go about repo root
)

func main() {
        fmt.Println(hello.Hello)
}
EOF

$ go mod tidy
go: finding module for package arch.seankhliao.com/hello.git
go: downloading arch.seankhliao.com/hello.git v0.0.0-20210329200831-02adee55f661
go: found arch.seankhliao.com/hello.git in arch.seankhliao.com/hello.git v0.0.0-20210329200831-02adee55f661

$ go run .
world
```
