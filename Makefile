.PHONY: dev build generate install image release profile bench test clean

CGO_ENABLED=0
COMMIT=$(shell git rev-parse --short HEAD)

all: dev

dev: build
	@./bitcask --version

build: clean generate
	@go build \
		-tags "netgo static_build" -installsuffix netgo \
		-ldflags "-w -X $(shell go list)/.Commit=$(COMMIT)" \
		./cmd/bitcask/...

generate:
	@go generate $(shell go list)/...

install: build
	@go install ./cmd/bitcask/...

image:
	@docker build -t prologic/bitcask .

release:
	@./tools/release.sh

profile: build
	@go test -cpuprofile cpu.prof -memprofile mem.prof -v -bench ./...

bench: build
	@go test -v -benchmem -bench=. ./...

test: build
	@go test -v -cover -coverprofile=coverage.txt -covermode=atomic -coverpkg=./... -race ./...

clean:
	@git clean -f -d -X
