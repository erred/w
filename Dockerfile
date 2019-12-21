FROM golang:alpine AS build

ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor
WORKDIR /app
COPY . .
RUN go build -o /bin/sitegen seankhliao.com/v10

FROM node:alpine

RUN apk --update --no-cache add imagemagick ca-certificates
RUN npm i -g firebase-tools
COPY --from=build /bin/sitegen /bin/sitegen

ENTRYPOINT ["/bin/sitegen"]
