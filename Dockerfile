FROM golang:alpine AS build

ENV CGO_ENABLED=0
WORKDIR /app
RUN apk --update --no-cache add imagemagick ca-certificates
COPY . .
RUN go build -mod=vendor -o /bin/sitegen ./sitegen

ENTRYPOINT ["/bin/sitegen"]
