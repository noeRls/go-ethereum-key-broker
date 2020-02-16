FROM golang:latest

WORKDIR /app

COPY keys.csv .

COPY get.sh .
RUN get.sh

COPY *.go .

CMD [ "go run ." ]