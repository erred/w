### _grpc_ for go

- install `protoc`, usually from package manager
- install `protoc-gen-go`: `go get google.golang.org/protobuf/cmd/protoc-gen-go`
  for message output
- install `protoc-gen-go-grpc`: `go get google.golang.org/grpc/cmd/protoc-gen-go-grpc`
  for service output

#### _proto_ file

```proto
syntax = "proto3";

package a.pkg;

option go_package = "github.com/example/pkg"

service Hello {
  rpc World(Msg) returns (Msg) {}
}

message Msg {}
```

#### _protoc_ invocation

the go generators currently take no useful options,
check by looking at generator's `main.go` and flags set there.

```sh
protoc \
  --proto_path=... \                    # search path for imports, repeatable
  --go_out=... \                        # go output directory
  --go_opt=paths=source_relative \      # place go files next to proto defs
  --go-grpc_out=... \                   # go-grpc output directory
  --go-grpc_opt=paths=source_relative \ # place go-grpc files next to proto defs
  *.proto                               # surce proto files
```

#### _linting_

- install `protoc-gen-lint`: `go get github.com/ckaznocha/protoc-gen-lint`

```sh
protoc --lint_out=sort_imports:. *.proto
```
