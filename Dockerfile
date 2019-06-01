FROM golang:alpine AS build

WORKDIR /app
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o app

FROM scratch

COPY --from=build /app/app /bin/app

ENTRYPOINT ["/bin/app"]
