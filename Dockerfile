FROM golang:rc-alpine AS build

ENV CGO_ENABLED=0
WORKDIR /workspace
COPY . .
RUN go run ./cmd/webrender && \
    go build -trimpath -ldflags='-s -w' -o /bin/serve ./cmd/serve

FROM scratch
COPY --from=build /bin/serve /bin/serve
ENTRYPOINT ["/bin/serve"]
