FROM golang:1.21-alpine as builder

COPY . /go/src

RUN cd src \
  && go mod download \
  && CGO_ENABLED=0 go build -v -a -ldflags "-s -w" \
  && mkdir dist \
  && cp -r fcct-service public dist/

FROM scratch

LABEL org.opencontainers.image.title="fcct-service" \
      org.opencontainers.image.version="0.19.0" \
      org.opencontainers.image.description="FCCT (butane) as a Service" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.source="https://github.com/zyclonite/fcct-service"

COPY --from=builder /go/src/dist /

EXPOSE 8080

ENTRYPOINT ["/fcct-service"]
