# Use uma imagem oficial do Go como a imagem base para a construção
FROM golang:1.24 

# Defina o diretório de trabalho dentro do container
WORKDIR /app

RUN apt-get update ; apt-get install coreutils nodejs -y

# Faça cache das dependências copiando go.mod e go.sum primeiro
COPY go.* ./
RUN go mod download
RUN go install github.com/air-verse/air@latest
# Copie o restante do código-fonte da aplicação
COPY . .

RUN chmod +x /app/cmd/start.sh

# Exponha a porta da aplicação
EXPOSE 9090

# Comando para executar a aplicação
# CMD ["air","--build.cmd","\"go build -o main main.go\"","--build.bin","\"./main\""]
CMD ["/go/bin/air"]