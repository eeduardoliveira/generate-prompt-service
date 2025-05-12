# Etapa 1: Build
FROM golang:1.23-alpine AS builder

# Define diretório de trabalho
WORKDIR /app

# Primeiro, copie apenas os arquivos de dependência
COPY go.mod go.sum ./

# Resolva as dependências (sem copiar o código ainda, melhora cache)
RUN go mod tidy && go mod download

COPY . .

# Compila a aplicação para produção
RUN go build -o generate-prompt-service ./main.go

# Etapa 2: Execução
FROM alpine:latest

WORKDIR /root/

# Copia binário da etapa de build
COPY --from=builder /app/generate-prompt-service .

CMD ["./generate-prompt-service"]