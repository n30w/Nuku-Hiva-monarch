# builder
FROM golang:1.19-alpine as builder

RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o andthensome .

# distributed image
FROM alpine:latest
COPY --from=builder /build/andthensome .

EXPOSE 80
EXPOSE 3306
EXPOSE 4000
ENTRYPOINT [ "./andthensome" ]
