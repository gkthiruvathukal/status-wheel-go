.PHONY: all build run test

all: build

build: build-server build-client

build-server:
	go build -o bin/server cmd/server/main.go

build-client:
	go build -o bin/client cmd/client/main.go

run: build
	./bin/server

test:
	go test ./pkg/status
