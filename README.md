[![Docker Pulls](https://badgen.net/docker/pulls/zyclonite/fcct-service)](https://hub.docker.com/r/zyclonite/fcct-service)

## fcct-service
FCCT as a Service

### build

`docker build -t zyclonite/fcct-service .`

### run

`docker run --name fcct-service -d -p 8080:8080 zyclonite/fcct-service`

### use

`curl -X POST --data-binary @test/fcos-config.yaml -H "Content-type: text/x-yaml" http://127.0.0.1:8080/api/v1/transpile?pretty=true&strict=false`

or open `http://127.0.0.1:8080/` in your browser for a simplistic ui

### demo

see [fcct.wsn.at](https://fcct.wsn.at)
