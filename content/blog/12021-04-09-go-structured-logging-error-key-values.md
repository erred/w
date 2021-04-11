---
title: go structured logging error key values
description: annotating your logging lines with key=value from errors
---

### _logging_

Logging, the more lightweight cousin to tracing everything,
and the formalized version of _printf_ debugging.

Anyway, if you believe in structured logging, in Go,
your logger interface should be [logr.Logger](https://pkg.go.dev/github.com/go-logr/logr#Logger):

```go
// github.com/go-logr/logr.Logger
type Logger interface {
        Enabled() bool
        Info(msg string, keysAndValues ...interface{})
        Error(err error, msg string, keysAndValues ...interface{})
        V(level int) Logger
        WithValues(keysAndValues ...interface{}) Logger
        WithName(name string) Logger
}
```

#### _error_

Now comes the problem:
you annotate errors with context, usually with something like:

```go
return fmt.Errorf("something went wrong foo=%q bar=%q: %w", hello, world, err)
```

But this doesn't fit well with structured logging,
you want the key-values to be independent,
not stuffed and escaped as a string in the error key.

So what can you do?

#### _ctxdata_

Since you already pass `context.Context` everywhere (you do, right?)
you might as well use it to smuggle data, ex using [ctxdata](https://pkg.go.dev/github.com/peterbourgon/ctxdata/v4)

```go
func ExampleCtxdata(l logr.Logger) {
        ctx := context.Background()
        ctx, data := ctxdata.New(ctx)

        err := ctxdataMid1(ctx)
        if err != nil {
                for _, kv := range data.GetAllSlice() {
                        l = l.WithValues(kv.Key, kv.Val)
                }
                l.Error(err, "some message", "example", "ctxerr")
        }
}

func ctxdataMid1(ctx context.Context) error {
        err := ctxdataMid2(ctx)
        if err != nil {
                data := ctxdata.From(ctx)
                data.Set("mid1", "midLevel")
                data.Set("repeated", "ctxdataMid1")
                return err
        }
        return nil
}

func ctxdataMid2(ctx context.Context) error {
        err := ctxdataDeep(ctx)
        if err != nil {
                return fmt.Errorf("unexpected: %w", err)
        }
        return nil
}

func ctxdataDeep(ctx context.Context) error {
        data := ctxdata.From(ctx)
        data.Set("deep", "value")
        data.Set("repeated", "ctxdataDeep")
        return errors.New("oops")
}

```

results in something like:
(using [klogr](https://pkg.go.dev/k8s.io/klog/v2/klogr) as the logger implementation).

```txt
E0411 11:42:37.052899   25625 ctxdata.go:21]  "msg"="some message" "error"="unexpected: oops" "deep"="value" "mid1"="midLevel" "repeated"="ctxdataMid1" "example"="ctxerr"
```

#### _custom_ error

Alternatively, you can use a custom error type that knows how to store errors:

```go
package sterr

import "errors"

type keyValuer interface {
        KeyValues() []interface{}
}

type structuredError struct {
        err error
        kvs []interface{}
}

func (e *structuredError) Error() string {
        return e.err.Error()
}

func (e *structuredError) KeyValues() []interface{} {
        // alternatively make this smarter and do unwrapping here
        // to support errors.As properly
        return e.kvs
}

func (e *structuredError) Unwrap() error {
        return e.err
}

func New(err error, keysAndValues ...interface{}) error {
        return &structuredError{
                err,
                keysAndValues,
        }
}

func FindKeyValues(err error) []interface{} {
        var kvs []interface{}
        if kverr, ok := err.(keyValuer); ok {
                kvs = append(kvs, kverr.KeyValues()...)
        }
        for err = errors.Unwrap(err); err != nil; err = errors.Unwrap(err) {
                if kverr, ok := err.(keyValuer); ok {
                        kvs = append(kvs, kverr.KeyValues()...)
                }
        }
        return kvs
}
```

and use it:

```go
func ExampleSterr(l logr.Logger) {
        err := sterrMid1()
        if err != nil {
                l.WithValues(sterr.FindKeyValues(err)...).Error(err, "some message", "example", "sterr")
        }
}

func sterrMid1() error {
        err := sterrMid2()
        if err != nil {
                return sterr.New(err, "mid1", "midlevel", "repeated", "sterrMid1")
        }
        return nil
}

func sterrMid2() error {
        err := sterrDeep()
        if err != nil {
                return fmt.Errorf("unexpected: %w", err)
        }
        return nil
}

func sterrDeep() error {
        return sterr.New(errors.New("oops"), "deep", "value", "repeated", "sterrDeep")
}
```

result, note the difference in repeated keys:

```txt
E0411 11:42:37.052779   25625 sterr.go:14]  "msg"="some message" "error"="unexpected: oops" "deep"="value" "mid1"="midlevel" "repeated"="sterrDeep" "example"="sterr"
```
