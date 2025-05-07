# Use uma imagem oficial do Go como a imagem base para a construção
FROM golang:1.24 AS builder

# Defina o diretório de trabalho dentro do container
WORKDIR /app

# Faça cache das dependências copiando go.mod e go.sum primeiro
COPY go.* ./
RUN go mod download

# Copie o restante do código-fonte da aplicação
COPY . .

# Compile a aplicação Go de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -o main .

# Use uma imagem base mínima para o container final
FROM alpine:latest

# Instale as dependências necessárias (se necessário)
RUN apk add --no-cache ca-certificates

# Defina o diretório de trabalho dentro do container
WORKDIR /app

# Copie o binário compilado da etapa de construção
COPY --from=builder /app/main .

COPY db/migrations/ db/migrations/
COPY .env .env

# Exponha a porta da aplicação
EXPOSE 9090

# Comando para executar a aplicação
CMD ["./main"]