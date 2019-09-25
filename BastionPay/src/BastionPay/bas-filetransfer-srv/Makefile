all: media

clean:
	@rm -rf ./api-article-86-x64-linux*

media:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v .
	@mv ./api-article ./api-article-86-x64-linux