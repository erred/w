FROM golang:alpine AS build

WORKDIR /app
RUN apk --update --no-cache add imagemagick ca-certificates
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o /bin/sitegen ./sitegen

ENTRYPOINT ["/bin/sitegen"]
