FROM golang:latest

WORKDIR /app
RUN mkdir keys_found

COPY keys_db ./keys_db/

COPY get.sh .
RUN sh get.sh

COPY *.go /app/

RUN go build -o go-ethereum-key-broker .

VOLUME [ "/app/keys_found" ]

CMD [ "./go-ethereum-key-broker" ]