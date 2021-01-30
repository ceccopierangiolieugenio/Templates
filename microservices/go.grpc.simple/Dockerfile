# Build
ARG folder_prj
FROM golang:1.15.7-alpine3.13 as builder

RUN apk add protoc git

ENV FOLDER_PRJ=${folder_prj}

WORKDIR /go/src/prj
RUN echo ${FOLDER_PRJ}/main.go yyy ${FOLDER_PRJ} xxx $FOLDER_PRJ
COPY ${FOLDER_PRJ}/main.go .
COPY testweb testweb

ENV GO111MODULE=on
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go \
           google.golang.org/grpc/cmd/protoc-gen-go-grpc

RUN protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           testweb/testweb.proto

RUN go mod init testprot.eugenio && go build main.go

# Create prod image
FROM alpine:3.13

EXPOSE 5000

COPY --from=builder /go/src/prj/main .

CMD ["./main"]