FROM golang:1.19-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app



# Download necessary Go modules
# COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Copy source code into /app
# COPY *.go ./
# COPY .env ./

# RUN go build -ldflags="-s -w" .
RUN go build -o main
EXPOSE 8080
CMD [ "/app/main" ]
