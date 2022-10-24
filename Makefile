BINARY_NAME=andthensome

build:
	go build .

run:
	clear
	go build -o ${BINARY_NAME} .
	./${BINARY_NAME}

# compile:
# 	echo "Compiling for every OS and Platform"
# 	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
# 	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
#	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go

all: run
