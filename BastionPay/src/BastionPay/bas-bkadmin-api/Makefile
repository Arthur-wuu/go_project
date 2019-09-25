all: media clean

clean:
    @rm -rf ./backend-serve*

media:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 godep go build -v .
	@mv ./backend-serve ./backend-serve-86-x64-linux
