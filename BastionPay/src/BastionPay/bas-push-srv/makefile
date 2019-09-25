all: clean media

clean:
	@rm -rf ./bas-push-srv

media:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v .




