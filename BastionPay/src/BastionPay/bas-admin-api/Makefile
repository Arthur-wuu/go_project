all: media

clean:
	@rm -rf ./basadmin-api-86-x64-linux-*

media:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 godep go build -v .
	@mv ./basadmin-api ./basadmin-api-86-x64
