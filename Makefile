BINARY_NAME=juice
BINARY_DIR=bin/juice-cli
CLI_DIR=cmd/juice-cli

build:
	go build -o ${BINARY_DIR}/${BINARY_NAME} ${CLI_DIR}/main.go # build for current platform

build_all:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_DIR}/${BINARY_NAME}-darwin ${CLI_DIR}/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_DIR}/${BINARY_NAME}-linux ${CLI_DIR}/main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_DIR}/${BINARY_NAME}-windows ${CLI_DIR}/main.go

run:
	./${BINARY_DIR}/${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_DIR}/${BINARY_NAME}-darwin
	rm ${BINARY_DIR}/${BINARY_NAME}-linux
	rm ${BINARY_DIR}/${BINARY_NAME}-windows

test:
	go test ./...

test_verbose:
	go test ./... -v

test_coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all

.PHONY: build
