[![Docker Pulls](https://badgen.net/docker/pulls/zyclonite/fcct-service)](https://hub.docker.com/r/zyclonite/fcct-service)

## fcct-service
The [Fedora CoreOS Config Transpiler](https://github.com/coreos/butane) (Butane) as a Service

An API to translate human readable Fedora CoreOS Configs (FCCs) into machine readable [Ignition](https://github.com/coreos/ignition) Configs. See the [getting
started](https://github.com/coreos/butane/blob/master/docs/getting-started.md) guide for how to use FCCT and the [configuration
specifications](https://github.com/coreos/butane/blob/master/docs/specs.md) for everything FCCs support.

### build

`docker build -t zyclonite/fcct-service .`

### run

`docker run --name fcct-service -d -p 8080:8080 zyclonite/fcct-service`

### use

`curl -X POST --data-binary @test/fcos-config.yaml -H "Content-type: text/x-yaml" http://127.0.0.1:8080/api/v1/transpile?pretty=true&strict=false`

or open `http://127.0.0.1:8080/` in your browser for a simplistic ui

### demo

see [fcct.wsn.at](https://fcct.wsn.at)
