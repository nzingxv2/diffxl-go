# File: Makefile
.PHONY: build lint test clean

BINARY=diffxl-go

build:
	go build -o ${BINARY} main.go

lint:
	golangci-lint run

test:
	go test -coverprofile=coverage.out ./...

clean:
	go clean
	rm -f ${BINARY}