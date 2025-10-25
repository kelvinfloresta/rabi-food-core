FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .

RUN useradd -m appuser
USER appuser

CMD ["./server"]
