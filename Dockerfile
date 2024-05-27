FROM golang:1.22-alpine
ENV GOOS=linux
WORKDIR /app

COPY .. .

WORKDIR  /app/

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

RUN chmod -R 777 .

WORKDIR /app/example

EXPOSE 2222

ENTRYPOINT CompileDaemon -build="go build -o example.o" -command="./example.o"