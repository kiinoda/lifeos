.DEFAULT_GOAL := help
.PHONY: build clean deploy

help:
	@echo "Available targets:"
	@echo "  build  - Build Go binary for Linux"
	@echo "  clean  - Remove built binary"
	@echo "  deploy - Clean, build, and deploy"
	@echo "  help   - Show this help"

build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap cmd/bootstrap/*.go

clean:
	rm -f ./bootstrap

deploy: clean build
	npx sls deploy --verbose
