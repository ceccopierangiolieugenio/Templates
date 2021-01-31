# Base building image
FROM golang:1.15.7-alpine3.13 AS builder

RUN apk add protoc git

ENV GO111MODULE=on
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go \
           google.golang.org/grpc/cmd/protoc-gen-go-grpc

WORKDIR /go/src/prj
COPY testweb testweb

RUN protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           testweb/testweb.proto

RUN go mod init testprot.eugenio


# Build the Backend
FROM builder AS builder_backend
COPY backend backend
RUN go build backend/main.go
CMD CGO_ENABLED=0  go test -v ./backend/...

# Make the Server
FROM builder AS builder_server
COPY server server
RUN go build server/main.go
CMD CGO_ENABLED=0  go test -v ./server/...


# Create prod backend image
FROM alpine:3.13 AS backend_image
COPY --from=builder_backend /go/src/prj/main .
CMD ["./main"]

# Create prod server image
FROM alpine:3.13 AS server_image
COPY --from=builder_server /go/src/prj/main .
CMD ["./main"]