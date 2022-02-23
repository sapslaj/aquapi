.PHONY: build clean deploy gomodgen

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o out/images functions/images/*.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o out/lucky functions/lucky/*.go

clean:
	rm -rf ./out ./vendor

deploy: clean build
	sls deploy --verbose
