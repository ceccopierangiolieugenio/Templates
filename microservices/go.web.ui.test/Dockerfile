# Make the Server
FROM golang:1.15.7-alpine3.13 AS builder_test_server

WORKDIR /go/src/prj
COPY server server
RUN go build server/main.go
CMD CGO_ENABLED=0  go test -v ./server/...


# Create prod server image
FROM alpine:3.13 AS server_image

COPY --from=builder_test_server /go/src/prj/main .
COPY client client
CMD ["./main"]