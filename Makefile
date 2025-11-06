binary_file::=./bin/aur-helper
module_path::=./cmd/aur/...
module_name::=$(shell sed -n 's/^module //p' go.mod)
git_version::=$(shell git describe --long --tags 2>/dev/null)
app_version::=$(shell if [ -n "${git_version}" ]; then echo "${git_version}" | sed -E 's/([^-]*)-g([0-9a-f]+)/\1+\2/'; else printf '0.0.0-%s+%s' "$(shell git rev-list --count HEAD)" "$(shell git rev-parse --short HEAD)"; fi)

.PHONY: help deps utils build lint lint-bin install clean

## help: print this help message
help:
	@echo 'Usage:'
	@sed -En 's/^##\s*//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/\t/'

## deps: install dependencies
deps:
	go mod tidy
	mkdir -p "$(shell dirname ${binary_file})"

## utils: install utils for subcommands
utils:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/deadcode@latest
	go install golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest

## build: build the application for production
build:
	CGO_ENABLED=0 GOEXPERIMENT=greenteagc,jsonv2 go build -ldflags="-s -w -X '${module_name}/internal/app.version=${app_version}'" -trimpath -a -installsuffix cgo -o "${binary_file}" "${module_path}"

## lint: run linters
lint:
	go vet ./...
	golangci-lint run ./...
	govulncheck -show=traces ./...
	deadcode -test ./...
	modernize -fix -test ./...

## lint-bin: run binary linters
lint-bin: build
	govulncheck -mode=binary -show=traces "${binary_file}"

DESTDIR ?= ./bin
## install: install the application
install:
	install -Dsm755 ${binary_file} "$(shell realpath $(DESTDIR))/$(shell basename ${binary_file})"

## clean: cleanup tasks
clean:
	rm -fR tmp bin "$(shell dirname ${binary_file})"
