#!/bin/bash
echo prepare $1...$2
echo GOPATH=$GOPATH
cd $GOPATH
git clone -b $2 https://github.com/BastionPay/$1.git src/BastionPay/$1
if [ $? -ne 0 ]
then
  echo error
  return
fi
cd src/BastionPay/$1
echo sync vendors...
govendor sync
if [ $? -ne 0 ]
then
  echo error
  return
fi
echo begin to make...
make
