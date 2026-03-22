# Estágio de Compilação
FROM golang:1.25.6-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o system-monitor ./cmd/system-monitor

# Estágio de Execução (Onde usamos o Alpine)
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/system-monitor .
EXPOSE 8080
CMD ["./system-monitor"]