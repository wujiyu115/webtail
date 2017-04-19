# WebTail

Realtime tail -f logs on the web

## Features
- [x] Use websocket to receive realtime log.
- [x] Support go-bindata.
- [ ] Support multi-files to display.

## Build

* go get github.com/jteeuwen/go-bindata/...
* go get golang.org/x/tools/cmd/goimports

* `./gobin.sh`
* `go build`

## Run

* `./webtail --log=/tmp/xxx.log`
./awebtail --log=/home/far/work/workspace/BnHServer/nohup.out
