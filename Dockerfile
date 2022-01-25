FROM golang:1.17-alpine as builder

COPY . /go/src

RUN cd src \
  && go mod download \
  && CGO_ENABLED=0 go build -v -a -ldflags "-s -w"

FROM scratch

LABEL org.opencontainers.image.title="fcct-service" \
      org.opencontainers.image.version="0.13.1" \
      org.opencontainers.image.description="FCCT (butane) as a Service" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.source="https://github.com/zyclonite/fcct-service"

COPY --from=builder /go/src/fcct-service /
COPY --from=builder /go/src/public /public/

EXPOSE 8080

ENTRYPOINT ["/fcct-service"]
