.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bootstrap cmd/bootstrap/*.go

clean:
	rm -f ./bootstrap

deploy: clean build
	npx sls deploy --verbose
