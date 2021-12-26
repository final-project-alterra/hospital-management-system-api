# syntax=docker/dockerfile:1

FROM golang:1.16-alpine3.14

WORKDIR /hms

COPY . .

RUN go mod download

RUN go build -o mainfile

EXPOSE 80/tcp

CMD ["./mainfile"]
