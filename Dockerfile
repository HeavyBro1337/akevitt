FROM golang:1.20-alpine
ENV GOOS=linux
WORKDIR /app

COPY . .

WORKDIR  /app/examples/iron-exalt

RUN mkdir data

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

RUN go build -o iron-exalt

RUN chmod 777 ./iron-exalt

EXPOSE 2222

ENTRYPOINT CompileDaemon -build="go build -o iron-exalt" -command="./iron-exalt"