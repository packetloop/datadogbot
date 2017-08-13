PROJECT_NAME := datadogbot
package = github.com/packetloop/$(PROJECT_NAME)

all: release

.PHONY: release
release: install
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/$(PROJECT_NAME)-linux-amd64 $(package)
	GOOS=darwin GOARCH=amd64 go build -o release/$(PROJECT_NAME)-darwin-amd64 $(package)

.PHONY: install
install:
	go get -v ./...
	go fmt ./...
	go vet ./...
	go test -race -v
	go build
	go build ./...
