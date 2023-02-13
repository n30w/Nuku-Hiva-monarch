BINARY_NAME=andthensome

build:
	go build .

run:
	clear
	go build -o ${BINARY_NAME} .
	./${BINARY_NAME}

populate:
	make run
	curl localhost:4000/populate

all: make run

test:
	clear
	go run .
	go test -v

lint:
	golangci-lint run
