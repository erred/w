---
description: using channels in place of other directives
title: go channel sync
---

### Synchronization

From [GopherCon 2018: Bryan C. Mills - Rethinking Classical Concurrency Patterns](https://www.youtube.com/watch?v=5zXAHh5tJqQ)

#### _exclusive_ access

```go
func ExampleChannelLock() {
        mu := make(chan int, 1)
        mu <- 0

        x := <- mu
        x = x * 2
        mu <- x
}()
```

#### _broadcast_

Alternatively start with `nil` inner channel for nonblocking start
(and check for nils).

```go
func ExampleBroadcast() {
        bc := make(chan chan struct{}, 1)
        bc <- make(chan struct{})

        // wait
        go func() {
                s := <-bc
                bc <- s
                select {
                case <-s:
                }
        }()

        // signal
        s := <-bc
        close(s)
        bc <- make(chan struct)
}
```

#### _limited_ concurrency

less useful when you may need to recursively add work

```go
func ExampleSemaphore() {
        sem := make(chan struct{}, limit)
        for _, task := range tasks {
                sem <- struct{}
                go func() {
                        work(task)
                        <-sem
                }()
        }

        // wait
        for n := limit; n > 0; n-- {
                sem <- struct{}
        }
}
```
