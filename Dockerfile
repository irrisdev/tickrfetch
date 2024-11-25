FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o tickrbot .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/tickrbot .

CMD ["./tickrbot"]