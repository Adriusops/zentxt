#multi-stage build
#stage 1: build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o zentxt ./cmd/zentxt/

#stage 2: run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/zentxt .

CMD ["./zentxt"]
