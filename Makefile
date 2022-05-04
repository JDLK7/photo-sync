clean:
	rm -rf dist

build-linux:
	mkdir -p dist/linux/armv7
	mkdir -p dist/linux/amd64
	GOOS=linux GOARCH=arm GOARM=7 go build -o dist/linux/armv7/photo-sync main.go
	GOOS=linux GOARCH=amd64 go build -o dist/linux/amd64/photo-sync main.go

build-darwin:
	mkdir -p dist/darwin/amd64
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/amd64/photo-sync main.go

build-all: build-linux build-darwin

dist: build-all
	tar czvf dist/photo-sync-linux-armv7.tar.gz dist/linux/armv7
	tar czvf dist/photo-sync-linux-amd64.tar.gz dist/linux/amd64
	tar czvf dist/photo-sync-darwin-amd64.tar.gz dist/darwin/amd64

.PHONY: clean build-linux build-darwin build-all dist
