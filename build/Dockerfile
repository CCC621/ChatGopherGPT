FROM golang:1.19.1-alpine

RUN apk update && apk add git

WORKDIR /go/src/cmd

RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/aws/aws-lambda-go/lambda@latest

