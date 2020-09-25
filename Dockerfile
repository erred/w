FROM golang:alpine AS build

WORKDIR /workspace
COPY go.mod .
COPY go.sum .
COPY cmd cmd
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /bin/http-server ./cmd/http-server

FROM us.gcr.io/com-seankhliao/webstyle:latest AS render

WORKDIR /workspace
COPY public public
COPY src src
RUN ["/bin/webrender"]

FROM scratch
COPY --from=build /bin/http-server /bin/http-server
COPY --from=render /workspace/public /var/public
ENTRYPOINT ["/bin/http-server", "-dir=/var/public"]
