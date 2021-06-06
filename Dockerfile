FROM golang:alpine AS build
ENV CGO_ENABLED=0
WORKDIR /workspace
COPY . .
RUN go run ./cmd/webrender -src content -dst static/root -gtm GTM-TLVN7D6
RUN go build -trimpath -ldflags='-s -w' -o /bin/w16 ./cmd/w16

FROM gcr.io/distroless/static
COPY --from=build /bin/w16 /bin/w16
ENTRYPOINT ["/bin/w16"]
