# Base building image
FROM golang:1.15.7-alpine3.13 AS builder

RUN apk add protoc git

ENV GO111MODULE=on
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go \
           google.golang.org/grpc/cmd/protoc-gen-go-grpc

WORKDIR /go/src/prj
COPY services/proto_pds/*.proto proto_pds/

RUN protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto_pds/pds.proto

RUN go mod init test.90poe/services



# Build the Backend
FROM builder AS builder_pds
COPY services/pds pds
RUN go build pds/cmd/backend/main.go

# Make the Server
FROM builder AS builder_client
COPY services/client client
RUN go build client/cmd/server/main.go



# Create prod pds image
FROM alpine:3.13 AS pds_image
COPY --from=builder_pds /go/src/prj/main .
CMD ["./main"]

# Create prod server image
FROM alpine:3.13 AS client_image
COPY --from=builder_client /go/src/prj/main .
CMD ["./main"]