FROM golang:alpine AS build

ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor
WORKDIR /app
RUN apk --update --no-cache add imagemagick ca-certificates
COPY . .
RUN go build -mod=vendor -o /bin/sitegen seankhliao.com/com-seankhliao/v7/sitegen

ENTRYPOINT ["/bin/sitegen"]
