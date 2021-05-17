FROM golang:1.16.3-buster

RUN mkdir /go/src/app
WORKDIR /go/src/app

ADD ./ /go/src/app
RUN apt update
RUN apt install -y make bash gcc g++ tzdata vim git libpq-dev
RUN go get
RUN go build main.go