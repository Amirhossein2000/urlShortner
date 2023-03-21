# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD ./cmd ./cmd
ADD ./internal ./internal

RUN go build -o /urlShortner ./cmd

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /urlShortner /urlShortner

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/urlShortner"]