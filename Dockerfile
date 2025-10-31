# Etapa 1: Build da aplicação
FROM golang:1.25.3-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/.env .
COPY --from=builder /app/login_code.html .

EXPOSE 6868

CMD ["./app"]
