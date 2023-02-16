# builder
FROM golang:1.19-alpine as builder

RUN mkdir /build
<<<<<<< HEAD
ADD ../../. /build
=======
ADD ../../ /build
>>>>>>> 0feba56bbad5bf5eeab343b92c98e3308c5b633e
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o andthensome ./cmd/andthensome/main.go

# distributed image
FROM alpine:latest
COPY --from=builder /build/andthensome .

EXPOSE 80
EXPOSE 3306
EXPOSE 4000
ENTRYPOINT [ "./andthensome" ]
