binary_file=./bin/aur-helper
module_path=./cmd/aur/...
module_name=$$(sed -n 's/^module //p' go.mod)
app_version=$$(set -o pipefail; git describe --long --tags 2>/dev/null | sed -r 's/([^-]*-g)/r\1/;s/-/./g' || printf "r%s.%s" "$$(git rev-list --count HEAD)" "$$(git rev-parse --short HEAD)")

.PHONY: help deps deps-lint build lint install clean

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## deps: install dependencies
deps:
	go mod tidy
	mkdir -p $$(dirname ${binary_file})

deps-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/deadcode@latest

## build: build the application for production
build: deps
	go build -ldflags="-s -w -X '${module_name}/app.version=${app_version}'" -a -installsuffix cgo -o ${binary_file} ${module_path}

lint: deps-lint
	golangci-lint run ./...
	govulncheck -show=traces ./...
	deadcode -test ./...

ifeq ($(DESTDIR),)
	DESTDIR := ./bin
endif

## install: install the binary file
install:
	install -Dsm755 ${binary_file} "$$(realpath $(DESTDIR))/$$(basename ${binary_file})"

## clean: cleanup tasks
clean:
	rm -fR $$(dirname ${binary_file})
	go clean -cache
