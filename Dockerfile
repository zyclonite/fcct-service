FROM golang:1.15-alpine as builder

COPY . /go/src

RUN cd src \
  && go mod download \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -ldflags "-s -w"

FROM scratch
LABEL version="0.10.0"
LABEL description="FCCT as a Service"

COPY --from=builder /go/src/fcct-service /
COPY --from=builder /go/src/public /public/

EXPOSE 8080

ENTRYPOINT ["/fcct-service"]
