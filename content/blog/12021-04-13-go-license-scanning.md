---
title: go license scanning
description: did i accidentally use any GPL things...
---

### _license_ scan

Open source is fun, but also kinda annoying.
See that thing? It's free and available to use,
except it has an incompatible license....

Anyway, find out if you're accidentally violating any today

#### _tools_

_go-licenses_ is probably your best bet right now.

##### _i_ tried

###### go-licenses

[go-licenses](https://github.com/google/go-licenses)

This one also tried to discover the url for the licenses,
drop stderr to ignore the errors from it not understanding vanity imports.

```sh
$ git clone https://github.com/google/go-licenses
$ cd go-licenses
$ go install

$ cd ~/w
$ go-licenses csv ./... 2>/dev/null
go.opentelemetry.io/contrib/instrumentation/host,Unknown,Apache-2.0
go.opentelemetry.io/otel/metric,Unknown,Apache-2.0
github.com/felixge/httpsnoop,https://github.com/felixge/httpsnoop/blob/master/LICENSE.txt,MIT
github.com/prometheus/common/internal/bitbucket.org/ww/goautoneg,https://github.com/prometheus/common/blob/master/internal/bitbucket.org/ww/goautoneg/README.txt,BSD-3-Clause
go.opentelemetry.io/otel/exporters/metric/prometheus,Unknown,Apache-2.0
github.com/golang/protobuf,https://github.com/golang/protobuf/blob/master/LICENSE,BSD-3-Clause
github.com/matttproud/golang_protobuf_extensions/pbutil,https://github.com/matttproud/golang_protobuf_extensions/blob/master/pbutil/LICENSE,Apache-2.0
github.com/prometheus/procfs,https://github.com/prometheus/procfs/blob/master/LICENSE,Apache-2.0
go.opentelemetry.io/otel,Unknown,Apache-2.0
github.com/cespare/xxhash/v2,https://github.com/cespare/xxhash/blob/master/v2/LICENSE.txt,MIT
github.com/prometheus/common,https://github.com/prometheus/common/blob/master/LICENSE,Apache-2.0
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp,Unknown,Apache-2.0
github.com/prometheus/client_golang/prometheus,https://github.com/prometheus/client_golang/blob/master/prometheus/LICENSE,Apache-2.0
github.com/beorn7/perks/quantile,https://github.com/beorn7/perks/blob/master/quantile/LICENSE,MIT
go.opentelemetry.io/otel/sdk/export/metric,Unknown,Apache-2.0
go.opentelemetry.io/otel/sdk,Unknown,Apache-2.0
github.com/yuin/goldmark-meta,https://github.com/yuin/goldmark-meta/blob/master/LICENSE,MIT
gopkg.in/yaml.v2,Unknown,Apache-2.0
github.com/shirou/gopsutil,https://github.com/shirou/gopsutil/blob/master/LICENSE,BSD-3-Clause
go.opentelemetry.io/contrib,Unknown,Apache-2.0
github.com/prometheus/client_model/go,https://github.com/prometheus/client_model/blob/master/go/LICENSE,Apache-2.0
golang.org/x/sys,Unknown,BSD-3-Clause
go.opentelemetry.io/otel/trace,Unknown,Apache-2.0
go.opentelemetry.io/contrib/instrumentation/runtime,Unknown,Apache-2.0
go.opentelemetry.io/otel/sdk/metric,Unknown,Apache-2.0
go.seankhliao.com/w/v15,Unknown,MIT
github.com/go-logr/logr,https://github.com/go-logr/logr/blob/master/LICENSE,Apache-2.0
google.golang.org/protobuf,Unknown,BSD-3-Clause
k8s.io/klog/v2,Unknown,Apache-2.0
github.com/yuin/goldmark,https://github.com/yuin/goldmark/blob/master/LICENSE,MIT
```

###### wwhrd

[wwhrd](https://github.com/frapposelli/wwhrd)

I think it's too fine grained? working on the package instead of the module/repo level.

```sh
$ git clone https://github.com/frapposelli/wwhrd
$ cd wwhrc
$ go install .

$ cd ~/w
$ go mod vendor
$ wwhrd list
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/trace
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/set
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark/extension
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/controller/time
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/instrumentation
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/encoding/protowire
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/genname
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/filetype
INFO[0001] Found License                                 license=MIT package=github.com/beorn7/perks/quantile
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/controller/basic
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/semconv
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/encoding/prototext
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/reflect/protoregistry
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/descfmt
INFO[0001] Found License                                 license=MIT package=github.com/go-ole/go-ole
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/procfs/internal/util
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/trace
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/fieldsort
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/internal/global
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/detrand
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/encoding/tag
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/encoding/messageset
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/aggregator/sum
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/internal/trace/parent
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/pragma
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/internal/baggage
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/common/model
INFO[0001] Found License                                 license=MIT package=github.com/StackExchange/wmi
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/procfs/internal/fs
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark-meta
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/golang/protobuf/ptypes/duration
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark/extension/ast
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/client_golang/prometheus/internal
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/shirou/gopsutil/cpu
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/common/internal/bitbucket.org/ww/goautoneg
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/shirou/gopsutil/internal/common
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount
INFO[0001] Found License                                 license=Apache-2.0 package=gopkg.in/yaml.v2
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/shirou/gopsutil/mem
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/errors
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark/text
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/runtime/protoimpl
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/golang/protobuf/ptypes/any
INFO[0001] Found License                                 license=BSD-3-Clause package=golang.org/x/sys/unix
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/internal
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/export/trace
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/encoding/text
INFO[0001] Found License                                 license=MIT package=github.com/felixge/httpsnoop
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/flags
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/propagation
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/shirou/gopsutil/process
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/attribute
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/metric/global
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/procfs
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/client_golang/prometheus/promhttp
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark/ast
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/metric
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/internal
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/unit
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/internal/metric
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark/util
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/metric/registry
INFO[0001] Found License                                 license=BSD-3-Clause package=golang.org/x/sys/internal/unsafeheader
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/golang/protobuf/ptypes/timestamp
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/version
INFO[0001] Found License                                 license=Apache-2.0 package=k8s.io/klog/v2/klogr
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/client_golang/prometheus
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/types/known/anypb
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/runtime/protoiface
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/export/metric
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/common/expfmt
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/fieldnum
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark/renderer
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/encoding/defval
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/selector/simple
INFO[0001] Found License                                 license=MIT package=github.com/cespare/xxhash/v2
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/aggregator/exact
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/strs
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/shirou/gopsutil/net
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark/parser
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/contrib
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/resource
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/contrib/instrumentation/runtime
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/filedesc
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/proto
INFO[0001] Found License                                 license=MIT package=github.com/go-ole/go-ole/oleutil
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/prometheus/client_model/go
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/aggregator
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/impl
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/aggregator/histogram
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/matttproud/golang_protobuf_extensions/pbutil
INFO[0001] Found License                                 license=BSD-3-Clause package=golang.org/x/sys/windows
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark/renderer/html
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/descopts
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/types/known/durationpb
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/golang/protobuf/proto
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/internal/mapsort
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/exporters/metric/prometheus
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/metric/processor/basic
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/internal/trace/noop
INFO[0001] Found License                                 license=Apache-2.0 package=k8s.io/klog/v2
INFO[0001] Found License                                 license=Apache-2.0 package=github.com/go-logr/logr
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/types/known/timestamppb
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/codes
INFO[0001] Found License                                 license=MIT package=github.com/yuin/goldmark
INFO[0001] Found License                                 license=BSD-3-Clause package=google.golang.org/protobuf/reflect/protoreflect
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/metric/number
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/contrib/instrumentation/host
INFO[0001] Found License                                 license=BSD-3-Clause package=github.com/golang/protobuf/ptypes
INFO[0001] Found License                                 license=Apache-2.0 package=go.opentelemetry.io/otel/sdk/export/metric/aggregation
```

###### golicense

[golicense](https://github.com/mitchellh/golicense)

This works on final binaries, so no test dependencies included...

also reaches out to github.

```sh
$ git clone https://github.com/mitchellh/golicense
$ cd golicense
$ go install .

$ cd ~/w
$ go build ./cmd/w
$ golicense -plain w
github.com/beorn7/perks                                       MIT License
github.com/prometheus/common                                  Apache License 2.0
github.com/cespare/xxhash                                     MIT License
github.com/golang/protobuf                                    BSD 3-Clause "New" or "Revised" License
github.com/shirou/gopsutil                                    BSD 3-Clause "New" or "Revised" License
github.com/prometheus/client_golang                           Apache License 2.0
github.com/prometheus/client_model                            Apache License 2.0
go.opentelemetry.io/contrib                                   Apache License 2.0
go.opentelemetry.io/otel                                      Apache License 2.0
go.opentelemetry.io/contrib/instrumentation/runtime           Apache License 2.0
github.com/felixge/httpsnoop                                  MIT License
github.com/go-logr/logr                                       Apache License 2.0
github.com/prometheus/procfs                                  Apache License 2.0
golang.org/x/sys                                              BSD 3-Clause "New" or "Revised" License
go.opentelemetry.io/otel/trace                                Apache License 2.0
google.golang.org/protobuf                                    BSD 3-Clause "New" or "Revised" License
go.opentelemetry.io/otel/metric                               Apache License 2.0
go.opentelemetry.io/contrib/instrumentation/host              Apache License 2.0
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp Apache License 2.0
go.opentelemetry.io/otel/exporters/metric/prometheus          Apache License 2.0
k8s.io/klog                                                   Apache License 2.0
github.com/matttproud/golang_protobuf_extensions              Apache License 2.0
go.opentelemetry.io/otel/sdk                                  Apache License 2.0
go.opentelemetry.io/otel/sdk/metric                           Apache License 2.0
go.opentelemetry.io/otel/sdk/export/metric                    Apache License 2.0
```

###### lc

[lc](https://github.com/boyter/lc)

only walks filesystem paths (works better if you vendor all your deps)

```sh
$ git clone https://github.com/boyter/lc
$ cd lc
$ ./generate_database.sh
$ go install .

$ cd ~/w
$ go mod vendor
$ lc
LICENSE
 likely licence; unable to identify
content/blog/12019-06-18-license-to-hack.md
 likely licence; unable to identify
content/blog/12021-04-13-go-license-scanning.md
 likely licence; unable to identify
vendor/github.com/StackExchange/wmi/LICENSE
 likely licence; unable to identify
vendor/github.com/beorn7/perks/LICENSE
 likely licence; unable to identify
vendor/github.com/cespare/xxhash/v2/LICENSE.txt
 likely licence; unable to identify
vendor/github.com/felixge/httpsnoop/LICENSE.txt
 likely licence; unable to identify
vendor/github.com/go-logr/logr/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/github.com/go-ole/go-ole/LICENSE
 likely licence; unable to identify
vendor/github.com/golang/protobuf/LICENSE
 Blended BSD-3-Clause 79.60701512785884
vendor/github.com/matttproud/golang_protobuf_extensions/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/github.com/prometheus/client_golang/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/github.com/prometheus/client_model/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/github.com/prometheus/common/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/github.com/prometheus/procfs/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/github.com/shirou/gopsutil/LICENSE
 Blended BSD-3-Clause 79.44269573919561
vendor/github.com/yuin/goldmark/LICENSE
 likely licence; unable to identify
vendor/github.com/yuin/goldmark-meta/LICENSE
 likely licence; unable to identify
vendor/go.opentelemetry.io/contrib/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/contrib/instrumentation/host/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/contrib/instrumentation/runtime/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/otel/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/otel/exporters/metric/prometheus/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/otel/metric/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/otel/sdk/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/otel/sdk/export/metric/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/otel/sdk/metric/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/go.opentelemetry.io/otel/trace/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/golang.org/x/sys/LICENSE
 Blended BSD-3-Clause 79.63331092357137
vendor/golang.org/x/sys/unix/syscall_bsd.go
 likely licence; unable to identify
vendor/golang.org/x/sys/unix/xattr_bsd.go
 likely licence; unable to identify
vendor/google.golang.org/protobuf/LICENSE
 likely licence; unable to identify
vendor/gopkg.in/yaml.v2/LICENSE
 Blended Apache-2.0 82.67973856209152
vendor/gopkg.in/yaml.v2/LICENSE.libyaml
 likely licence; unable to identify
vendor/k8s.io/klog/v2/LICENSE
 likely licence; unable to identify
```

###### go-license-detector

[go-license-detector](https://github.com/go-enry/go-license-detector)

Doesn't seem to find anything from my dependencies?

```sh
$ git clone https://github.com/go-enry/go-license-detector
$ cd go-license-detector
$ go install ./cmd/license-detector

$ cd ~/w
$ license-detector .
.
        95%        MIT
        82%        MIT-0
```

##### _other_ tools

###### snyk

[snyk](https://snyk.io/)

apparently it can do this? but I think you have to link your source code to their online service...

###### licensed

[licensed](https://github.com/github/licensed)

don't know how to run ruby, also works on file paths only (vendor)
