# WebTail

Realtime tail -f logs on the web

## Features
- [x] Use websocket to receive realtime log.
- [x] Support go-bindata.
- [x] Support error filter

## Build

* go get github.com/jteeuwen/go-bindata/...
* go get golang.org/x/tools/cmd/goimports

* `./build.sh`

## Run

* `./webtail --log=/tmp/xxx.log`
