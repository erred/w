---
title: go private modules and gitlab subgroups
description: why do so many people get tripped up on this
---

### _subgroups_

Say you use [GitLab](https://gitlab.com/)
for hosting code and you use their subgroup/team thingy to organize people (or projects).
How do you get it work with Go modules (privately!)

#### _dependency_

Say the thing you want to import is at

```txt
gitlab.com/my-team/my-subteam/my-subsubteam/repo-a
```

witht the following code

```go
-- go.mod --
module gitlab.com/my-team/my-subteam/my-subsubteam/repo-a

go 1.16

-- a.go --
package a

var A = "A"
```

#### _use_ it

So, like all private code,
set the `GOPRIVATE` environment variable so `go` doesn't reach out to a proxy

```sh
export GOPRIVATE=gitlab.com/my-team
# or
go env -w GOPRIVATE=gitlab.com/my-team
```

Setup git to clone using ssh,
`go` clones using https but you need a way to authenticate it.

```gitconfig
[url "git@gitlab.com:"]
    insteadOf = "https://gitlab.com/"
```

And obtain an access token from gitlab and put it in `~/.netrc`,
so the `go get` can find the correct repo root

_important:_ the access token needs `read_api`, not just `read_repository`

```netrc
machine gitlab.com
login seankhliao
password _FREdNJyBnFwZDn9Gj48
```

And you're all set

```sh
go mod tidy
```

#### _troubleshooting_

##### _gitconfig_

```sh
» go get
go: gitlab.com/testgroup-395/foo@v0.0.0-20210429153539-e2d639ad297e: invalid version: git fetch -f origin refs/heads/*:refs/heads/* refs/tags/*:refs/tags/* in /tmp/gomodcache.3BvK/cache/vcs/35cd6107a6e6f51d91c6a96d27cc113a029c4667736e65b94c2cd7c4dcf6d9ab: exit status 128:
        fatal: could not read Username for 'https://gitlab.com': terminal prompts disabled
```

You haven't setup noninteractive git clones, see above section for `gitconfig`.
You should be able to `git clone https://gitlab.com/you/private/repo` without any further input (importantly over https)

##### _netrc_

```sh
» go get
go: gitlab.com/testgroup-395/subgroup-a/bar@v0.0.0-20210429153603-7e4a416f18f5: invalid version: git fetch -f origin refs/heads/*:refs/heads/* refs/tags/*:refs/tags/* in /tmp/gomodcache.GqcM/cache/vcs/15b3f321a38509dd7662d8b9d4b7ad5dccf082b4d75d9646fe2fc1fc4fe59365: exit status 128:
        client_global_hostkeys_private_confirm: server gave bad signature for RSA key 0
        remote:
        remote: ========================================================================
        remote:
        remote: The project you were looking for could not be found or you don't have permission to view it.
        remote:
        remote: ========================================================================
        remote:
        fatal: Could not read from remote repository.

        Please make sure you have the correct access rights
        and the repository exists.
```

You haven't set `~/.netrc` or the token in `~/.netrc` doesn't have the correct scopes.
Check the output of:

- `go get -v` to see if it identified the correct repo root
- `cat $GOMODCACHE/cache/vcs/xxxxxx...xxx/config` (from above output) to see if it identified the correct repo root
- `curl --netrc https://gitlab.com/path/to/repo?go-get=1` to see if gitlab responds correctly / your token has sufficient privileges

##### _goprivate_

```sh
» go get
go: gitlab.com/testgroup-395/foo@v0.0.0-20210429153539-e2d639ad297e: verifying go.mod: gitlab.com/testgroup-395/foo@v0.0.0-20210429153539-e2d639ad297e/go.mod: reading https://sum.golang.org/lookup/gitlab.com/testgroup-395/foo@v0.0.0-20210429153539-e2d639ad297e: 410 Gone
        server response:
        not found: gitlab.com/testgroup-395/foo@v0.0.0-20210429153539-e2d639ad297e: invalid version: git fetch -f origin refs/heads/*:refs/heads/* refs/tags/*:refs/tags/* in /tmp/gopath/pkg/mod/cache/vcs/35cd6107a6e6f51d91c6a96d27cc113a029c4667736e65b94c2cd7c4dcf6d9ab: exit status 128:
                fatal: could not read Username for 'https://gitlab.com': terminal prompts disabled
```

You haven't set `GOPRIVATE` and `go` is asking the public sumdb for the checksum,
which it can't access, (private code, remember?).
