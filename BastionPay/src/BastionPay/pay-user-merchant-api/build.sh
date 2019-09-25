#! /bin/bash

export LOGCAT=true

set -e

APP_NAME=basadmin-api

BIN_APP=./$APP_NAME

if [ -f "$BIN_APP" ];then
    echo "delete $BIN_APP ok"
    rm -f $BIN_APP
fi

godep go build -v .

$BIN_APP -c=config.yml -logtostderr

#gowatch -o ./kyc -p ./src/gopkg.exa.center/kyc -args='-c=config/local/config.yaml,start'
