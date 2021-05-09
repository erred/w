---
title: go mod vendor non go things
description: how to get your other things included into your vendor dir
---

### _go_ mod vendor

`go mod vendor` is fairly limited in scope,
it will only include what's needed from your dependencies to build your main program.
So your dependencies tests and non go files all get trimmed out.
But sometimes, you want to keep some of those extra files,
like html or protobuf files.

Now with `embed` in 1.16, there's a way to make those files required as part of the build.

#### _old_ situation

To recap, usually only `.go` files needed to build your main program
are included in `vendor`.

```txt
» exa -T
.
├── testrepo-dependency
│  ├── go.mod
│  ├── hello.proto
│  ├── index.html
│  ├── v.go
│  └── v_test.go
└── testrepo-main
   ├── go.mod
   ├── main.go
   └── vendor
      ├── go.seankhliao.com
      │  └── testrepo-dependency
      │     └── v.go
      └── modules.txt
```

#### _embed_ hack

We can force files to be needed as part of a build by embedding them,
so in the above example in `v.go`:


```go
import "embed"

//go:embed *.html *.proto *.go
var _ embed.FS
```

`//go:embed *` won't work because `.git/` exists
but can't be embedded as it's not in the scope of a module.
This results in:

```txt
» exa -T
.
├── testrepo-418
│  ├── go.mod
│  ├── hello.proto
│  ├── index.html
│  ├── v.go
│  └── v_test.go
└── testrepo-419
   ├── go.mod
   ├── main.go
   ├── testrepo-419
   └── vendor
      ├── go.seankhliao.com
      │  └── testrepo-418
      │     ├── hello.proto
      │     ├── index.html
      │     ├── v.go
      │     └── v_test.go
      └── modules.txt
```
