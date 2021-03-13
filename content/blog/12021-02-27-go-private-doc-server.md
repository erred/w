---
description: docs like pkg.go.dev for your private go code
title: go private doc server
---

### _docs_

One nice thing about Go code is that docs are consistent,
for all public code, go to [pkg.go.dev](https://pkg.go.dev),
type in the import path (or append directly to the url),
and have docs.
No wondering how to navigate the funky webpage,
or even trying to find one,
if it's public, it's there.

#### _old_ way

But what about your private code?
Previously, you'd run [godoc](https://golang.org/x/tools/cmd/godoc)
in your `GOPATH`, and it'd serve everything in there,
plus the standard library.
But it's `GOPATH` and has no concept of versioning.

There's also [gddo](https://go.googlesource.com/gddo/),
for the extra scale needed to run [godoc.org](https://godoc.org).

#### _pkgsite_

The new way is to run a private version of _pkg.go.dev_,
the project is called [pkgsite](https://go.googlesource.com/pkgsite/).

##### _run_

You need a database to do this,
specifically postgres.
Also migrations,
and maybe redis, and a worker?
why is this so hard.
How do you even get data into the database?

_TODO:_ figure out how to do this, `docker-compose up` is not enough

##### _direct_

If you have an internal proxy serving all your code,
pkgsite can serve directly from that,
bypassing the need for a database.

_warning:_
pretty much the only thing that works
is viewing the docs for packages that you know the exact url for.
Search,
navigating to subdirectories in modules,
source code links all don't work.

```sh
git clone https://go.googlesource.com/pkgsite
cd pkgsite
go build golang.org/x/pkgsite/cmd/frontend
./frontend -bypass_license_check -direct_proxy
```

##### _local_ directory

There's also the experimental command to more or less do what _godoc_ did:
serve local directories without versions:

```sh
git clone https://go.googlesource.com/pkgsite
cd pkgsite
go build golang.org/x/pkgsite/cmd/pksite
./pkgsite /path/to/module1,path/to/module2/,...
```
