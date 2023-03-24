#start building this image from this source (OS, etc...)
FROM golang:1.19-alpine3.17

COPY . /go/src/go-docker

WORKDIR /go/src/go-docker

COPY go.mod ./

RUN go mod download

COPY . .

RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
RUN ["go", "install", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -polling -log-prefix=false -build="go build main.go" -command="./main"