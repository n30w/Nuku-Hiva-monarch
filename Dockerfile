FROM golang:1.19-alpine

RUN mkdir /app
ADD . /app

# Sets command context for subsequent commands
WORKDIR /app

# Download necessary Go modules
RUN go mod download

# RUN go build -ldflags="-s -w" .
RUN go build -o main
EXPOSE 8080
CMD [ "/app/main" ]
