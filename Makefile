binary_file::=./tmp/aur-helper
module_path::=./cmd/aur/...
module_name::=$(shell sed -n 's/^module //p' go.mod)
app_version::=$(shell git_output=$(shell git describe --long --tags 2>/dev/null); if [ "${?}" = 0 ]; then printf '%s' "${git_output}" | sed -r 's/([^-]*-g)/r\1/;s/-/./g'; else printf '0.0.0+r%s.%s' "$(shell git rev-list --count HEAD)" "$(shell git rev-parse --short HEAD)"; fi)

.PHONY: help deps utils build lint install clean

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## deps: install dependencies
deps:
	go mod tidy
	mkdir -p "$(shell dirname ${binary_file})"

## utils: install utils for subcommands
utils:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/deadcode@latest

## build: build the application for production
build:
	go env -w CGO_ENABLED=0
	go build -ldflags="-s -w -X '${module_name}/app.version=${app_version}'" -a -trimpath -installsuffix cgo -o ${binary_file} ${module_path}

## lint: run linters
lint: utils
	golangci-lint run ./...
	govulncheck -show=traces ./...
	deadcode -test ./...

DESTDIR ?= ./bin
## install: install the application
install:
	install -Dsm755 ${binary_file} "$(shell realpath $(DESTDIR))/$(shell basename ${binary_file})"

## clean: cleanup tasks
clean:
	rm -fR "$(shell dirname ${binary_file})"
	go clean -cache
