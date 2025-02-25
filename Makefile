default: check

## Main Commands

build: fmt clean test build-all package

clean: clean-pkg clean-bin tidy
	go clean -i -cache -testcache

## Supporting Commands

tidy:
	go mod tidy

fmt: tidy
	trunk fmt

fmt-all: tidy
	trunk fmt --all

check: fmt
	trunk check

check-all: fmt-all
	trunk check --all

test:
	go test ./lib/...

clean-bin:
	rm -f bin/*

clean-pkg:
	rm -f pkg/*

update: upgrade
upgrade: tidy
	go get -u
	trunk upgrade

## Build sub-commands

build-all: build-win build-linux

build-linux:
	GOOS=linux GOARCH=amd64 go build -o "bin/pdt" ./main.go

build-win:
	GOOS=windows GOARCH=amd64 go build -o "bin/pdt.exe" ./main.go

## Package sub-commands

package:
	zip -r "pkg/pdt.zip" "bin/pdt"
	zip -r "pkg/pdt.exe.zip" "bin/pdt.exe"

## Git Hooks

pre-commit: clean check test
	git add go.mod go.sum
