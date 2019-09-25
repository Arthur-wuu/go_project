all: clean media

clean:
	@rm -rf ./pay-user-merchant-api

media:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Meta=build_date=`date -u '+%Y-%m-%d_%H:%M:%S'`&build_revision=`git rev-parse HEAD`" -v .