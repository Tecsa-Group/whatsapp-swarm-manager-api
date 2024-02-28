# Use a imagem oficial do Golang
FROM golang:1.22.0-alpine AS builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Crie o diretório .ssh
RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh

# Copie os arquivos id_rsa e id_rsa.pub do diretório ssh-keys local para o contêiner
COPY ssh-keys/id_rsa /root/.ssh/id_rsa
COPY ssh-keys/id_rsa.pub /root/.ssh/id_rsa.pub

# Defina as permissões apropriadas para os arquivos
RUN chmod 600 /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa.pub

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
