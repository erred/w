---
description: building go apps in containers faster
title: benchmark kaniko go
---

### _kaniko_

[kaniko](https://github.com/GoogleContainerTools/kaniko)
a container builder from google, intended to be run as a container.
kaniko wants to be root but does not need privileges,
different from others.

docker with buildkit was not tested becasue
clean builds with both `--cache-from` and `--build-arg BUILDKIT_INLINE_CACHE=1` specified don't work.

#### _results_

_tldr_:
multistage builds with vendored deps have the fasted build times,
produce the slimmest images
but at the cost of bloating up the source code repo with checked in deps.

`--use-new-run` did not appear to have appreciable impact for multi-vendored.

| stages  | setup                        | uncached | cached     | dep update | cached `--use-new-run` | size     |
| ------- | ---------------------------- | -------- | ---------- | ---------- | ---------------------- | -------- |
| single  | basic                        | 171s     | 169s       | 163s       | 113s                   | 261.64MB |
| single  | go.mod/sum + go mod download | 201s     | build fail | 208s       | 194s                   | 317.15MB |
| single  | go.mod/sum + go mod vendor   | 181s     | build fail | 184s       | 160s                   | 264.95MB |
| single  | vendored                     | 176s     | 117s       | 112s       | _93s_                  | 149.77MB |
| multi   | basic                        | 126s     | 125s       | 113s       | 118s                   | _5.06MB_ |
| multi   | go.mod/sum + go mod download | 162s     | build fail | 141s       | 136s                   | _5.06MB_ |
| multi   | go.mod/sum + go mod vendor   | 140s     | build fail | 130s       | 123s                   | _5.06MB_ |
| _multi_ | _vendored_                   | _111s_   | _96s_      | _85s_      | _94s_                  | _5.06MB_ |

build fails were kaniko complaing about unlinkat / device busy,
did not appear to be transient failures.

```txt
error building image: error building stage: failed to execute command: extracting fs from image: removing whiteout .wh.workspace: unlinkat //workspace: device or resource busy
```

#### _setup_

Repo: GitHub
CI: Google Cloud Build, triggered by app

#### _Dockerfile_

Dockerfiles for the different tests,
uncomment as necessary

All tests used tidied go.mod/go.sum.
The vendored tests used the same setup as basic but with a checked in `vendor/` directory.

```dockerfile
FROM golang:alpine AS build

WORKDIR /workspace
RUN apk add --update --no-cache ca-certificates

# go.mod/sum + go mod download
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# go.mod/sum + go mod vendor
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /bin/app

# multistage
# FROM scratch AS app
# COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# COPY --from=build /bin/app /bin/

ENTRYPOINT [ "/bin/app" ]
```

#### _cloudbuild.yaml_

cloudbuild config for testing.
only the `_IMG` variable differed,
also `--use-new-run` for the test.

```yaml
substitutions:
  _IMG: test-multi-basic
  _REG: us.gcr.io
tags:
  - $SHORT_SHA
  - $COMMIT_SHA
steps:
  - id: build-push
    name: gcr.io/kaniko-project/executor:latest
    args:
      - -c=.
      - -f=Dockerfile
      - -d=$_REG/$PROJECT_ID/$_IMG:latest
      - -d=$_REG/$PROJECT_ID/$_IMG:$SHORT_SHA
      - --cache=true
      - --reproducible
      # - --use-new-run
```

#### _code_

The following code has a total of 112 dependencies (!!)
weighing in at 17M (vendored).

go.mod:

```gomod
module go.seankhliao.com/testrepo-176

go 1.16

require (
        go.seankhliao.com/usvc v0.8.8
        golang.org/x/net v0.0.0-20200822124328-c89045814202
)
```

main.go:

```go
package main

import (
        "flag"
        "os"

        "go.seankhliao.com/usvc"
        "golang.org/x/net/context"
)

var (
        cacheBust = 0
)

func main() {
        usvc.Exec(context.Background(), &Server{}, os.Args)
}

type Server struct{}

func (s *Server) Flags(fs *flag.FlagSet)                        {}
func (s *Server) Setup(ctx context.Context, u *usvc.USVC) error { return nil }
```
