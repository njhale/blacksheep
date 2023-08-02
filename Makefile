.PHONY: build-linux 

BUILD_LINUX := GOOS=linux GOARCH=amd64

build-linux:
	@$(BUILD_LINUX) go build -o blacksheep
