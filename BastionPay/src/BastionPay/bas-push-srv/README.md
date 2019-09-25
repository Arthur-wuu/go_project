1.install govendor<br/>
> brew install govendor

2.set GOPATH<br/>
> export GOPATH=~/gopath

3.goto $GOPATH<br/>
> cd $GOPATH

4.download<br/>
> git clone github.com/BastionPay/bas-push-srv.git src/BastionPay/bas-push-srv

5.goto source
> cd src/BastionPay/bas-push-srv

6.switch your branch<br/>
> git checkout master/dev

7.sync dependency
> govendor sync

8.make build<br/>
> make

use build.sh to build<br/>
> export GOPATH=~/gopath

build master<br/>
> ./build.sh bas-push-srv master

build dev<br/>
> ./build.sh bas-push-srv dev
