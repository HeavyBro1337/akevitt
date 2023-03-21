FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN mkdir data

RUN go mod download


RUN go build -o /akevitt

EXPOSE 2222

CMD [ "/akevitt" ]