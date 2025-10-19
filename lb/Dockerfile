FROM golang:1.24 AS build

COPY ./src /app/src

# Create output directory
RUN mkdir -p build

WORKDIR /app/src

RUN rm go.mod
RUN go mod init lb
RUN go build -o build/lb core/load_balancer.go

CMD ["./build/lb"]
