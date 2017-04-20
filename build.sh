#!/usr/bin/env bash

echo "go-bindata template files..."
go-bindata data/...
go build
./webtail --log=/home/far/work/workspace/BnHServer/nohup.out