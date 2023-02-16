BINARY_NAME=andthensome

# Creates a distributable binary.
bin:
	go build -o ${BINARY_NAME} ./cmd/andthensome/main.go

# Builds and runs the binary.
run:
	make bin
	./${BINARY_NAME}

all: make run

test:
	go test -v ./internal/**

lint:
	golangci-lint run

vendoring:
	go mod tidy
	go mod vendor
