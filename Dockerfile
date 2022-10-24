# builder
FROM golang:1.19-alpine as builder

RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o andthensome .
# RUN go build -o main .

# Clean image
FROM alpine:3.14
COPY --from=builder /build/andthensome .
COPY --from=builder /build/.env .

EXPOSE 80
EXPOSE 3306
EXPOSE 4000
# ENV PORT 3306
ENTRYPOINT [ "./andthensome" ]
# CMD [ "/app/main" ]
# CMD [ "./andthensome" ]