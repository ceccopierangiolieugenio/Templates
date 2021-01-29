# Build
# FROM golang:1.15.7 as builder
FROM golang:1.15.7-alpine3.13 as builder

WORKDIR /go/src/webserver
COPY app.go .

RUN go build app.go

# Create prod image
# FROM golang:1.15.7-alpine3.13
FROM alpine:3.13

EXPOSE 5000

COPY --from=builder /go/src/webserver/app .

CMD ["./app"]