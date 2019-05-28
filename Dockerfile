FROM golang:alpine AS build

WORKDIR /app
COPY site-builder .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o app

FROM scratch

WORKDIR /app
COPY --from=build /app/app .

WORKDIR /workspace
ENTRYPOINT ["/app/app"]
