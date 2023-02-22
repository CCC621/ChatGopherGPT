FROM golang:1.19.1-alpine

RUN apk update && apk add git

WORKDIR /go/src/cmd/bot

RUN go mod init ChatGopherGPT
RUN go get github.com/bwmarrin/discordgo@latest 
RUN go get github.com/aws/aws-lambda-go/lambda 
RUN go get github.com/huggingface/go-gpt3








