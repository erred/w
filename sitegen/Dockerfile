FROM golang:alpine AS build

WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o app


FROM alpine:latest

RUN apk --update --no-cache add imagemagick ca-certificates
COPY --from=build /app/app /bin/sitegen

ENTRYPOINT ["/bin/sitegen"]
