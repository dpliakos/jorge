FROM golang:1.18.3-alpine

RUN mkdir -p /go/src/github.com/dpliakos/jorge
WORKDIR /go/src/github.com/dpliakos/jorge

RUN apk add make

COPY ./main.go  .
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./go.mod .
COPY ./go.sum .
COPY ./LICENSE .
COPY ./Makefile .

RUN go get
RUN make build
RUN make install 

RUN mkdir /root/projectRoot
WORKDIR /root/projectRoot
