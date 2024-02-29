# Use a imagem oficial do Golang
FROM golang:1.22.0-alpine AS builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Crie o diretório .ssh
RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh

# Define um argumento para a chave privada
ARG PRIVATE_KEY

ARG POSTGRES_HOST
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_DB
ARG POSTGRES_PORT
ARG CLOUDFLARE_API_TOKEN
ARG CLOUDFLARE_ZONE_ID
ARG EVOLUTION_APIKEY

RUN echo "POSTGRES_HOST=$POSTGRES_HOST" > .env && \
    echo "POSTGRES_USER=$POSTGRES_USER" >> .env && \
    echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD" >> .env && \
    echo "POSTGRES_DB=$POSTGRES_DB" >> .env && \
    echo "POSTGRES_PORT=$POSTGRES_PORT" >> .env && \
    echo "CLOUDFLARE_API_TOKEN=$CLOUDFLARE_API_TOKEN" >> .env && \
    echo "CLOUDFLARE_ZONE_ID=$CLOUDFLARE_ZONE_ID" >> .env && \
    echo "EVOLUTION_APIKEY=$EVOLUTION_APIKEY" >> .env

# Escreve o valor do argumento no arquivo id_rsa
RUN echo "$PRIVATE_KEY" > /root/.ssh/id_rsa

# Defina as permissões apropriadas para os arquivos
RUN chmod 600 /root/.ssh/id_rsa



# Copie o arquivo go.mod e go.sum para baixar as dependências
COPY go.mod go.sum ./

# Baixe as dependências do módulo
RUN go mod download

# Copie o restante do código-fonte da aplicação
COPY . .

RUN chmod +x ./deploy_stack.sh

# Compile o binário da aplicação
RUN go build -o main .

EXPOSE 5000

# Comando para executar a aplicação
CMD ["./main"]
