# usvc

A repo for usvc

[![License](https://img.shields.io/github/license/seankhliao/usvc.svg?style=flat-square)](LICENSE)
![Version](https://img.shields.io/github/v/tag/seankhliao/usvc?sort=semver&style=flat-square)
[![pkg.go.dev](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/go.seankhliao.com/usvc)

```go
package main

import (
        "flag"
        "net/http"
        "os"

        "go.seankhliao.com/usvc"
)

type Server struct {
        svc *usvc.ServerSimple

        msg string
}

func NewServer(args []string) *Server {
        fs := flag.NewFlagSet(args[0], flag.ExitOnError)
        s := &Server{
                svc: usvc.NewServiceSimple(usvc.NewConfig(fs)),
        }
        fs.StringVar(&s.msg, "msg", "a message", "a message to the world")
        fs.Parse(args[1:])
        s.svc.Mux.HandleFunc("/", s.messenger)
        return s
}

func (s *Server) messenger(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(s.msg))
        s.svc.Log.Info().Msg("spread the message")
}

func main() {
        usvc.Run(usvc.SignalContext(), NewServer(os.Args).svc)
}
```
