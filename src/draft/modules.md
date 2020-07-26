---
description: practical go modules
title: go modules
---

### _go modules_

Go modules are a collection of packages which share a
_single unit of versioning_.

design goal: given the same set of initial dependencies
always select the same set of all dependencies.
(in contrast with SAT solvers which try to select the latest version of all dependencies)

#### _basics_

check in `go.mod` and `go.sum`

##### _semver_

Go modules depends on modules following semver as defined by
[semver.org](https://semver.org/).

`vX.Y.Z-a+b`

- `X`: major version, increment on breaking change
- `Y`: minor version, increment on new feature
- `Z`: patch version, increment on bugfix
- `a`: optional prerelease specifier, ordered alphabetically
- `b`: optional build specifier, ignored in version selection

`v0` and `v1` share the same namespace and allow seamless upgrading

##### _version_ selection

example: the following modules / versions are available

```txt
module-a v0.1.0 v0.2.0 v0.3.0 v1.0.0 v1.1.0-alpha.1
module-b v1.4.0 v1.5.0 v1.6.0 v2.0.0 v2.1.0
module-c v0.7.0 v0.8.0 v0.9.0
```

###### _new_ dependency

if unspecified, select the latest released version with the current major version

_example_:

- adding a new dependency `module-a`: `v1.0.0` is selected
- adding a new dependency `module-b`: `v1.6.0` is selected
- adding a new dependency `module-b/v2`: `v2.1.0` is selected

###### _indirect_ / existing dependency

_example_:

- `module-a` requires `module-c` @ `v0.7.0`
- `module-b` requires `module-c` @ `v0.8.0`
- `v0.8.0` of `module-c` is selected

##### _proxy_

By default the go toolchain will use a public proxy `https://proxy.golang.org`.
This offers faster querying and download speeds.

- `go env GOPROXY` check current settings
- once things are in the proxy they cannot be deleted
- practice with prerelease tags, you don't want to mess up the name you want to use
- moving tags or rewriting history will likely invalidate checksums resulting in errors
- the proxy has caches, bypass with `GOPROXY=direct` (fetches directly from source)
- private modules are not available through the proxy, set `GOPRIVATE`
- if you're developing multiple modules with constantly updated versions you may also want to set `GOPRIVATE` to bypass the cache
- RedHat patched the toolchain to bypass the default proxy

##### _multi module_ repositories

- Don't
- modules have no parent-child / hierarchical relationships
- just don't. seriously. don't.

#### _go.mod_

##### _create_ a module

in `go.mod`:

```go.mod
module example.com/module
```

the first part of the name _SHOULD_ have a dot in the name,
preferably domain like so `go get` can work,
see _vanity import paths_.
Only `test/` and `example/` are reserved to not need a dot.

- `go mod init example.com/module`
- `go mod init example.com/module/v2`
- `go mod init github.com/seankhliao/testrepo`
- `go mod init github.com/seankhliao/testrepo/v3`
- `go mod init test/module`

##### _go_ directive

in `go.mod`:

```go.mod
go 1.X
```

Where `11 <= X <= current version`.
Lower versions guarantee the availability of features at that version,
but restrict the availability of newer versions.
Go does not automatically download older/newer toolchains,
but the current one will report a mismatch if it fails to build.

##### _require_ directive

- `require example.com/module-a v1.2.3`
- `require example.com/module-a v2.3.4 // +incompatible`
- `require example.com/module-a v1.2.3-20200231-sdndsjcn`

##### _replace_ directive

are only applied in the main module (the `go build` is being run from)

- `replace example.com/module-a v1.2.3 => example.com/module-a v1.0.0`
- `replace example.com/module-a => exmaple.com/module-a v1.0.0`
- `replace example.com/module-a => example2.com/module-a v1.2.3`
- `replace example.com/module-a => ../module-a`

##### _exclude_ directive

- `exclude example.com v1.2.3`

#### _go.sum_

`go.sum` is not used for module selected / dependency resolution.
It is purely for security / integrity.
Do not worry about what's inside

#### _dependency_ management

##### _add_ a dependency

Add the dependency in your code with `import "example.com/module"`,
and run `go run .`, `go build`, `go mod tidy` to automatically update `go.mod`.
The latest version will be selected and recorded (if not already there).

- `go get example.com/module@latest`
- `go get example.com/module@v1.2.3`
- `go get example.com/module@master` (or any other branch)
- `go get example.com/module@sj4jdn3` (or any commit)

##### _update_ a dependency

- `go get -u=patch example.com/module`
- `go get -u example.com/module`
- `go get -u all`
- `go get -u ./...`

##### _remove_ a dependency

Remove all imports in your code,
run `go mod tidy`.

#### _vendor_

Vendoring is not necessary to guarantee the version of modules used,
but you may want it to colocate your dependencies with your code,
for auditing reasons, or for guaranteeing the availability of dependencies.

- create or sync: `go mod vendor`
- go will rewrite `vendor/` contents, do not modify
- go >= 1.14 will use `vendor/` automatically if the directory is available
- go &lt;= 1.13 needs the `-mod=vendor` flag

#### _release_ new versions

##### _patch_ versions

tag a new release

##### _minor_ versions

tag a new release

##### _major_ versions

###### v0 -> v1

tag a new release

###### v2+

update `go.mod` to new major version,

- `module example.com/module` => `module example.com/module/v2`
- `module example.com/module/v2` => `module example.com/module/v3`

update all internal imports to use the new import path
(unless you want to depend on an older version of yourself)

- `import "example.com/module/package"` => `import "example.com/module/v2/package"`
- `import "example.com/module/v2/package"` => `import "example.com/module/v3/package"`

#### _backwards_ compatibility

Sometimes you need to support legacy versions of go.

##### _do nothing_

update `go.mod` and imports as above, tag a new release.

- go modules users will import through `example.com/module/v2`
- gopath users will import through `example.com/module`
- gopath users will not be able to easily remain on previous versions

##### release branch

cut a new branch for every release
(major / minor / patch depends on how much granularity you want to support).

- go modules users will import through `example.com/module/v2`
- gopath users will import through `example.com/module`
- gopath users can easily get and use a previous release

##### go1 branch

similar to the above approaches except the `go1` branch has higher priority than the default branch

##### major subdirectory

Keep all existing code.
Create a `vX` subdirectory where `X` == major version, ex `v2`.
`go mod init example.com/module/vX` in the directory and copy all the other code in.

- go modules users will import through `example.com/module/v2`
- gopath users will import through `example.com/module/v2`
- bloats your repo

#### _querying_ modules

- `go mod why`
- `go mod graph`
- `go list -m -f {{ ... }}`
