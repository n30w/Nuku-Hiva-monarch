BINARY_NAME=andthensome

# Creates a distributable binary.
bin:
	go build -o ${BINARY_NAME} ./cmd/andthensome/main.go

# Builds and runs the binary.
run:
	make bin
	./${BINARY_NAME}

populate:
	make run
	curl localhost:4000/populate

all: make run

test:
	go test -v ./internal/**

lint:
	golangci-lint run
