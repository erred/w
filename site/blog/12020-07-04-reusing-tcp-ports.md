---
description: how to reuse ports for TCP in Go
title: reusing tcp ports
---

### _TCP_ port reuse

Usually, when you make outgoing TCP connections,
you don't care about the source address,
and you can just get a free ephemeral one.

Sometimes you do,
and you want all your connections to reuse the same source
ip:port, after all, TCP connections only need to be a unique 4 tuple.

Sometimes you want even more,
you want to accept incoming connections on a port,
and also make outgoing connections from the same port.

#### _Basics_

The long way round, with syscalls

##### _Normal_ dial

- `socket`: to get a new file descriptor, fd
- `bind`: to associate fd with local address
- `connect`: to associate fd with remote address

```go
func dial(la *syscall.SockaddrInet4) net.Conn {
        syscall.ForkLock.Lock()
        fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
        if err != nil {
                panic(err)
        }
        syscall.ForkLock.Unlock()

        err = syscall.Bind(fd, la)
        if err != nil {
                panic(err)
        }

        ra := &syscall.SockaddrInet4{Port: 9999}
        copy(ra.Addr[:], net.IPv4(127, 0, 0, 1))

        err := syscall.Connect(fd, ra)
        if err != nil {
                panic(err)
        }

        conn, err := net.FileConn(os.NewFile(uintptr(fd), ""))
        if err != nil {
                panic(err)
        }

        conn.Write([]byte("foo bar"))
}
```

##### _Normal_ listen

- `socket`: to get a new file descriptor
- `bind`: to associate fd with local address
- `listen`: to allow incoming connections
- `accept`: (?)

```go
func listen(sa *syscall.SockaddrInet4) {
        syscall.ForkLock.Lock()
        fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
        if err != nil {
                panic(err)
        }
        syscall.ForkLock.Unlock()

        err = syscall.Bind(fd, la)
        if err != nil {
                panic(err)
        }

        err = syscall.Listen(fd, 10)
        if err != nil {
                panic(err)
        }

        l, err := net.FileListener(os.NewFile(uintptr(fd), ""))
        if err != nil {
                panic(err)
        }

        for {
                conn, err := l.Accept()
                if err != nil {
                        panic(err)
                }
                conn.Write([]byte("hello world"))
        }
}
```

##### _Options_

`setsockopt`: modify how the socket will behave, use before `bind`

there are 2 options we care about:

- `unix.SO_REUSEADDR`: allows bind to reuse addresses for outgoing requests, also skips wait for timeout
- `unix.SO_REUSEPORT`: allows multiple sockets to use the same address

```go
err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, unix.SO_REUSEADDR|unix.SO_REUSEPORT, 1)
if err != nil {
  panic(err)
}
```

#### _net_

Thankfully `net` is flexible enough that
we don't need to do everything with syscalls

##### _dialer_

```go
func dial() {
        la, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
        if err != nil {
                panic(err)
        }

        d := &net.Dialer{
                LocalAddr: la,
                Control: func(network, address string, c syscall.RawConn) error {
                        var err error
                        c.Control(func(fd uintptr) {
                                err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEADDR|unix.SO_REUSEPORT, 1)
                        })
                        return err
                },
        }

        conn, err := d.Dial("tcp", "127.0.0.1:10000")
        if err != nil {
                panic(err)
        }

        conn.Write([]byte("foo bar"))
}
```

##### _listener_

```go
func listen() {
        lc := &net.ListenConfig{
                Control: func(network, address string, c syscall.RawConn) error {
                        var err error
                        c.Control(func(fd uintptr) {
                                err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEADDR|unix.SO_REUSEPORT, 1)
                        })
                        return err
                },
        }

        l, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:9999")
        if err != nil {
                panic(err)
        }

        for {
                conn, err := l.Accept()
                if err != nil {
                        panic(err)
                }
                conn.Write([]byte("hello world"))
        }
}
```
