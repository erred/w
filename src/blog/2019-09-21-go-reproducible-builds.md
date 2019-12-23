--- title
go reproducible builds
--- description
reproducible builds in Go
--- main


1. use **Go1.13+**
2. `CGO_ENABLED=0`: haven't had the time to figure out otherwise
3. `go mod vendor`
   - modules and proxies should help but this is the only way to make sure nothing changes out from under you, _have all your sources_
4. `go build -mod=vendor -trimpath`
   - trims the filesystem of the build environment from the resulting binary
