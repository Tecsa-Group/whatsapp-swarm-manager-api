# Use a imagem oficial do Golang
FROM golang:1.22.0-alpine AS builder

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

# Segunda etapa - Imagem final
FROM alpine:latest

# Atualize o índice de pacotes do Alpine Linux
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Copie o binário compilado da primeira etapa
COPY --from=builder /app/main /usr/local/bin/main

# Exponha a porta em que a aplicação vai rodar
EXPOSE 8088

# Comando para executar a aplicação
CMD ["main"]
