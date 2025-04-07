FROM golang:1.23.8-alpine AS builder

COPY ./src /go/src/

WORKDIR /go/src

RUN go mod tidy
RUN go build -o ./bin main.go

FROM ubuntu:bionic
RUN apt-get update && apt-get install -y docker.io
COPY --from=builder /go/src /go/src
WORKDIR /go/src

CMD ./bin

EXPOSE 9999
