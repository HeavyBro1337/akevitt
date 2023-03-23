FROM golang:1.20-alpine
ENV GOOS=linux
WORKDIR /app

COPY . .

WORKDIR  /app/samples/iron-exalt

RUN mkdir data

RUN go mod download


RUN go build -o iron-exalt

RUN chmod 777 ./iron-exalt

EXPOSE 2222

CMD [ "./iron-exalt" ]