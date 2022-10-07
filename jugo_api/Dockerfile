## build 

FROM golang:1.18 AS build

WORKDIR /jugo_api

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /jugo_app 

## deploy stage 2

FROM gcr.io/distroless/base-debian10

WORKDIR / 

COPY --from=build /jugo_app /jugo_app

EXPOSE 4000

USER nonroot:nonroot

ENTRYPOINT [ "/jugo_app" ]

