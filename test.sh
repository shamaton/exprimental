#!/bin/sh

DIR=$(cd $(dirname $0) && pwd)
cd ${DIR}

GOPATH=${DIR}
cd src/experimental

go test -v .
echo ""
sleep 3
go test -bench . -benchmem