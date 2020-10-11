FROM golang:alpine AS build

WORKDIR /workspace
COPY go.mod .
COPY go.sum .
COPY vendor vendor
COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /bin/ ./...
COPY public public
COPY site site
RUN ["/bin/webrender"]

FROM scratch
COPY --from=build /bin/com-seankhliao /bin/com-seankhliao
COPY --from=build /workspace/public /var/public
ENTRYPOINT ["/bin/com-seankhliao", "-dir=/var/public"]
