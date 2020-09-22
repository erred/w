FROM golang:alpine AS build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /bin/stream-server ./cmd/stream-server

FROM scratch

# sqlite?
COPY --from=build /etc/services /etc/services
COPY --from=build /etc/protocols /etc/protocols

COPY --from=build /bin/stream-server /bin/

ENTRYPOINT ["/bin/stream-server"]
