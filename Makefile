BINARY  := soji
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -s -w -X github.com/YOU/soji/internal/version.Version=$(VERSION)

.PHONY: build run test lint tidy clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/soji

run: build
	./bin/$(BINARY)

test:
	go test ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

clean:
	rm -rf bin
