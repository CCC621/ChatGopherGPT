FROM golang:1.19.1-alpine

RUN apk update && apk add git

WORKDIR /go/src/cmd/bot

RUN go mod init ChatGopherGPT
RUN go get github.com/bwmarrin/discordgo@latest 
RUN go get github.com/aws/aws-lambda-go/lambda 
RUN go get github.com/huggingface/go-gpt3







<<<<<<< Updated upstream
=======
# go mod init
# go get github.com/bwmarrin/discordgo@latest
# go get github.com/PullRequestInc/go-gpt3@latest
# go install golang.org/x/tools/cmd/goimports@latest
# go install github.com/aws/aws-lambda-go/lambda@latest
>>>>>>> Stashed changes

CMD ["go", "run", "main.go"]
