# syntax=docker/dockerfile:1

FROM golang:1.20-alpine

WORKDIR /src

COPY  . .

RUN go build -o /bin/blacksheep

EXPOSE 8080

ENTRYPOINT ["/bin/blacksheep"]
