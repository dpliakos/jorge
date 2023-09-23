FROM golang:1.18.3-alpine as build

WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download

COPY . .
RUN go build

FROM alpine:3.17.5 as final

USER 1000
COPY --from=build /app/jorge /jorge

WORKDIR /projectRoot
ENTRYPOINT ["/jorge"]
