FROM golang:alpine AS build
ARG CGO_ENABLED=0
ARG GOPROXY=https://proxy.golang.org,direct
ARG GOMODCACHE=/go/pkg/mod
ARG GOCACHE=/root/.cache/go-build
ARG GTM=GTM-TLVN7D6
WORKDIR /workspace
COPY . .
RUN go run ./cmd/webrender -src content -dst static/root -gtm ${GTM}
RUN go build -trimpath -ldflags='-s -w' -o /usr/local/bin/ ./cmd/...

FROM gcr.io/distroless/static AS singlepage
COPY --from=build /usr/local/bin/singlepage /bin/singlepage
ENTRYPOINT ["/bin/singlepage"]

FROM gcr.io/distroless/static AS w16
COPY --from=build /usr/local/bin/w16 /bin/w16
ENTRYPOINT ["/bin/w16"]
