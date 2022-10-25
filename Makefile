BINARY_NAME=andthensome

build:
	go build .

run:
	clear
	go build -o ${BINARY_NAME} .
	./${BINARY_NAME}


all: run
