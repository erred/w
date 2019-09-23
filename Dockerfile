FROM golang:alpine AS build

ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY . .
RUN go build -mod=vendor -o app seankhliao.com/com-seankhliao/v7/sitegen


FROM alpine:latest

RUN apk --update --no-cache add imagemagick ca-certificates
COPY --from=build /app/app /bin/sitegen

ENTRYPOINT ["/bin/sitegen"]
