FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cocktaildle_backend cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/cocktaildle_backend .

EXPOSE 8080

CMD ["./cocktaildle_backend"]
