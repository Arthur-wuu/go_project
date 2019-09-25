#! /bin/bash

set -e

export LOGCAT=true

export GOPATH=${GOPATH}:/golang


APP_NAME=api-article

BIN_APP=./bin/${APP_NAME}

if [ -f "$BIN_APP" ];then
    echo "delete ${BIN_APP} ok"
    rm -f ${BIN_APP}
fi
#
#
#godep go build -v -o ${BIN_APP} .
#
#${BIN_APP} -c=conf/local/config.yml start

gowatch -o ./bin/${APP_NAME} -p . -args='-c=conf/local/config.yml,start'