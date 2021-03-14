FROM golang:alpine AS build

ENV CGO_ENABLED=0
WORKDIR /workspace
COPY . .
RUN go run ./cmd/webrender -src content -dst internal/static/root
RUN go build -trimpath -ldflags='-s -w' ./cmd/w

FROM scratch
COPY --from=build /workspace/w /bin/w
ENTRYPOINT ["/bin/w"]
