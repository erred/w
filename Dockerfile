FROM us.gcr.io/com-seankhliao/webstyle:latest AS builder

WORKDIR /workspace
COPY . .
RUN ["/bin/webrender"]

FROM us.gcr.io/com-seankhliao/http-server:latest
COPY --from=builder /workspace/public /var/public
ENTRYPOINT ["/bin/http-server", "-dir=/var/public"]
