---
description: iterators in go
title: go iterators
---

### _iterators_

#### _language_

the language is special, it gets to use `range`.

##### _maps_

```go
func maps() {
        m := make(map[string]string)

        for key := range m {
                _ = key
        }

        for key, value := range m {
                _, _ = key, value
        }
}
```

##### _slices_

```go
func slices() {
        s := make([]string, 0)

        for index := range s {
                _ = index
        }

        for index, value := range s {
                _, _ = index, value
        }
}
```

##### _channels_

```go
func channels() {
        c := make(chan string)

        for value := range c {
                _ = value
        }
}
```

#### _standard-ish_

these are the usually seen ones

##### _no_ error

No errors,
no need to handle any iteration errors.

```go
func noError() {
        s := bufio.NewScanner(nil)
        for s.Scan() {
                _ = s.Text()
        }
}
```

##### _deferred_ error

You keep any errors internally,
stopping iteration on error,
and providing a way to check iteration errors later.
Users may forget to check errors.

```go
func deferredError() {
        r := sql.Rows{}

        for r.Next() {
                _ = r.Scan()
        }
        err := r.Err()
        if err != nil {
                // handle iteration error
                _ = err
        }
}
```

##### _immediate_ error

You return any iteration errors immediately,
forcing users to to use a 3 valued for.
Technically the first iteration can also be in the for loop,
but then you need to predeclare the error and value var to
get the right scope to handle errors.

```go
func immediateError() {
        r := tar.NewReader()

        h, err := r.Next()
        for ; err != nil; h, err = r.Next() {
                _ = h
        }
        if err != nil {
                // handle error
                _ = err
        }

        // alternative
        var err error
        var h *tar.Header
        for h, err = r.Next(); err != nil; h, err = r.Next() {
                _ = h
        }
        // handle error
        // most likely needs to handle the "no more results" error differently
        _ = err
}
```

#### _exotic_

You want to be cool and use `range` too.

##### _channel_ abuse

maps and slices need the entire thing present at the start of iteration,
but channels....

This does not expose a way to return errors,
other than it being part of the iteration value.
Also this will leak a goroutine and channel if iteration ever stops early
(goroutine holds reference to channel because it is blocked on send, can't be GC'ed).

The leak could be handled by instead returning a read-write channel
and catching a write-on-closed channel panic
or by passing a done chan to the iterator to signal an end to iteration.

No wonder nobody does this,
and people recommend not returning channels as part of public APIs.

```go
type X struct{}
func (x X) Iterator() <-chan string {
        c := make(chan string)
        go func(){
                for {
                        // TODO: conditionally beakout
                        c <- generateValue()
                }
                close(c)
        }()
        return c
}

func exoticChannelIterator() {
        var x X

        for value := range x.Iterator() {

        }
}
```
