VERSION=v0.1.0

all: build

build-darwin-amd64:
	mkdir -p ./bin && GOOS=darwin GOARCH=amd64 go build -o bin/diffcsv-darwin-amd64-$(VERSION)

build-linux-amd64:
	mkdir -p ./bin && GOOS=linux GOARCH=amd64 go build -o bin/diffcsv-linux-amd64-$(VERSION)

build-windows-amd64:
	mkdir -p ./bin && GOOS=windows GOARCH=amd64 go build -o bin/diffcsv-windows-amd64-$(VERSION).exe

clean:
	rm -fr ./bin

build: build-darwin-amd64 build-linux-amd64 build-windows-amd64

.PHONY: build all