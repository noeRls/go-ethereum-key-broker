FROM golang:latest

WORKDIR /app
RUN mkdir keys_found

COPY keys.csv .

COPY get.sh .
RUN sh get.sh

COPY *.go /app/

CMD [ "go", "run", "." ]