FROM golang:1.24 AS build
COPY . /app/backends
WORKDIR /app/backends
RUN mkdir -p build
RUN go build -o build/backend backend.go