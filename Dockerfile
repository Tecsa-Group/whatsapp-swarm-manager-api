# Use a imagem oficial do Golang
FROM golang:1.22.0-alpine AS builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o arquivo go.mod e go.sum para baixar as dependências
COPY go.mod go.sum ./

# Baixe as dependências do módulo
RUN go mod download

# Copie o restante do código-fonte da aplicação
COPY . .

# Compile o binário da aplicação
RUN go build -o main .

EXPOSE 5000

# Comando para executar a aplicação
CMD ["./main"]
