# Стадия сборки
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/

# Финальная стадия - минимальный образ без Go toolchain
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /server /app/server

EXPOSE 8080

ENTRYPOINT ["/app/server"]