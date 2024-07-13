BINARY_NAME=ypsilon14

build-macos-intel:
	GOARCH=amd64 GOOS=darwin go build -o bin/$(BINARY_NAME)-macos-intel

build-macos-arm:
	GOARCH=arm64 GOOS=darwin go build -o bin/$(BINARY_NAME)-macos-arm

build-linux:
	GOARCH=amd64 GOOS=linux go build -o bin/$(BINARY_NAME)-linux

build-windows:
	GOARCH=amd64 GOOS=windows go build -o bin/$(BINARY_NAME).exe

build-all: build-macos-intel build-macos-arm build-linux build-windows
