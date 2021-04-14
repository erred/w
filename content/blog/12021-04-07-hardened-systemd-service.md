---
title: hardened systemd service
description: like containers, but less convenient
---

### _systemd_

Why do you need containers?
"just run it as a systemd service with super-long-list
of options to lock it down", they said.

Anyway, I run Arch Linux, the unofficial testing ground of systemd,
so we have all the latest bugs as well.
As always, RTFM:
[system.unit](https://man.archlinux.org/man/systemd.unit.5.en)
[systemd.service](https://man.archlinux.org/man/systemd.service.5)
[systemd.exec](https://man.archlinux.org/man/systemd.exec.5.en)
[systemd.resource-control](https://man.archlinux.org/man/systemd.resource-control.5.en)

#### _app_

So, what are we containing?
A very generic web app exersizing the network, filesystem, and secrets.

_important:_ build as a static binary (`CGO_ENABLED=0`)

`main.go`:

```go
package main

import (
        "flag"
        "io"
        "log"
        "net/http"
        "os"
        "path/filepath"
)

func main() {
        var configPath, dataPath, certFile, keyFile string
        flag.StringVar(&configPath, "config", "", "path to config file")
        flag.StringVar(&dataPath, "data", "", "path to data dir")
        flag.StringVar(&certFile, "cert", "", "path to cert file")
        flag.StringVar(&keyFile, "key", "", "path to key file")
        flag.Parse()

        log.Printf("starting with config=%q data=%q", configPath, dataPath)
        b, err := os.ReadFile(configPath)
        if err != nil {
                log.Fatal("read config", err)
        }
        log.Println("config:\n", string(b))

        statefile := filepath.Join(dataPath, "statefile")

        http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
                if r.Method == http.MethodGet {
                        b, err := os.ReadFile(statefile)
                        if err != nil {
                                rw.WriteHeader(http.StatusInternalServerError)
                                rw.Write([]byte(err.Error()))
                                return
                        }
                        rw.Write(b)
                } else if r.Method == http.MethodPost {
                        defer r.Body.Close()
                        b, err := io.ReadAll(r.Body)
                        if err != nil {
                                rw.WriteHeader(http.StatusInternalServerError)
                                rw.Write([]byte(err.Error()))
                        }
                        err = os.WriteFile(statefile, b, 0o644)
                        if err != nil {
                                rw.WriteHeader(http.StatusInternalServerError)
                                rw.Write([]byte(err.Error()))
                        }
                        rw.Write([]byte("ok"))
                }
        })

        log.Println(http.ListenAndServeTLS(":8080", certFile, keyFile, nil))
}
```

#### _service_

So what do you ship to your target system?
A binary and a systemd service file of course.
No need for extra user, directory, etc.. setup,
but you'll want to provide your config files, and secrets.

`foobar.service`:

```systemd
[Unit]
Description=A foo bar app
Documentation=https://seankhliao.com/blog/12021-04-07-hardened-systemd-service/
Requires=network-online.target
After=network-online.target

[Service]
Type=simple
# path within the filesystem namespace
ExecStart=/bin/foobar -config /etc/foobar/conf.yaml -data /var/lib/foobar -cert ${CREDENTIALS_DIRECTORY}/cert -key ${CREDENTIALS_DIRECTORY}/key
RestartSec=60s
Restart=always

WorkingDirectory=/
# path on host system to be root of filesystem namespace,
# created by RuntimeDirectory=
# There's nothing in this namespace, including no libc,
# which is why statically linked binaries are important
# otherwise the error is a very opaque 203/EXEC File not found
RootDirectory=/run/foobar
ProtectProc=noaccess
ProcSubset=pid
# mount our executable from host into namespace
BindReadOnlyPaths=/usr/local/bin/foobar:/bin/foobar

DynamicUser=true

CapabilityBoundingSet=
AmbientCapabilities=

NoNewPrivileges=true

UMask=0022

ProtectSystem=strict
ProtectHome=true
RuntimeDirectory=foobar
StateDirectory=foobar
# CacheDirectory=
# LogDirectory=
ConfigurationDirectory=foobar
PrivateTmp=true
PrivateDevices=true
PrivateIPC=true
PrivateUsers=true
ProtectHostname=true
ProtectClock=true
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectKernelLogs=true
ProtectControlGroups=true
RestrictAddressFamilies=AF_INET AF_INET6
RestrictNamespaces=true
LockPersonality=true
MemoryDenyWriteExecute=true
RestrictRealtime=true
RestrictSUIDSGID=true
RemoveIPC=true
PrivateMounts=true

# I think this is everything most Go processes would need
# maybe you can trim it down more?
# SystemCallFilter=@basic-io @file-system @io-event @network-io @sync
SystemCallErrorNumber=EPERM
SystemCallArchitectures=native

# unset all env
Environment=

# Fancy credential passing
LoadCredential=cert:/etc/certs/example.com.pem
LoadCredential=key:/etc/certs/example.com-key.pem

IPAddressAllow=any
IPAddressDeny=
DevicePolicy=closed
DeviceAllow=


[Install]
WantedBy=mult-user.target
```

##### _update_

So if you make outgoing network calls, you'll want DNS and ca-certs.
Also might want tzdata? but Go can embed that these days.

```systemd
BindReadOnlyPaths=/usr/bin/feed-agg:/bin/feed-agg \
    /etc/feed-agg:/etc/feed-agg:rbind \
    /etc/resolv.conf:/etc/resolv.conf \
    /etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt
```
