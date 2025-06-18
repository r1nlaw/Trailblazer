FROM golang:1.23.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/server/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main ./
COPY schema/* ./schema/
COPY configs/* ./configs/
COPY images/* ./images/
COPY .env ./
COPY smtp-cert.pem ./

RUN chmod +x /app/main

EXPOSE 8080

RUN ls ./configs
CMD ["./main", "-c", "/app/configs/config.yml"]