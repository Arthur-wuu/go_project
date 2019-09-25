#! /bin/bash
echo prepare $1...$2
echo GOPATH=$GOPATH
cd $GOPATH
if [ ! -d "src/BastionPay/$1" ]; then
  echo first cloning...
  git clone -b $2 https://github.com/BastionPay/$1.git src/BastionPay/$1
  if [ $? -ne 0 ]
  then
    echo error
    return
  fi
else
  echo pulling...
  cd src/BastionPay/$1
  git checkout $2
  git pull
fi

cd $GOPATH/src/BastionPay/$1
echo sync vendors...
govendor sync
if [ $? -ne 0 ]
then
  echo error
  return
fi
echo begin to make...
make