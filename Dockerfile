FROM golang:1.14-alpine as builder

COPY . /go/src

RUN cd src \
  && go mod download \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -ldflags "-s -w"

FROM scratch
LABEL version="0.5.0"
LABEL description="FCCT as a Service"

COPY --from=builder /go/src/fcct-service /

EXPOSE 8080

ENTRYPOINT ["/fcct-service"]
