# builder
FROM golang:1.19-alpine as builder
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go mod download
RUN go build -ldflags="-s -w" -o andthensome .
# RUN go build -o main .

# Clean image
FROM alpine:3.14
COPY --from=builder /build/andthensome .
EXPOSE 8080
ENTRYPOINT [ "./andthensome" ]
# CMD [ "/app/main" ]
