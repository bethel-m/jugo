## build 

FROM golang:1.18 AS build

WORKDIR /jugo_api

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /jugo_app 



EXPOSE 4000

ENTRYPOINT [ "/jugo_app" ]

