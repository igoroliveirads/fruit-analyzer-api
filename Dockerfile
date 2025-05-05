# Estágio de build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Estágio final
FROM alpine:latest

WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/main .

# Expõe a porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"] 