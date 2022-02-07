.PHONY: build clean deploy gomodgen

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o out/images functions/images/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o out/lucky functions/lucky/main.go

clean:
	rm -rf ./out ./vendor

deploy: clean build
	sls deploy --verbose
