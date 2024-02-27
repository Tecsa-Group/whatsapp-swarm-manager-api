#!/bin/bash
echo "entrou porra"
set -e

# Verifica se um argumento foi fornecido
if [ "$#" -ne 2 ]; then
    echo "Número inválido de argumentos"
    echo "Uso: $0 <SERVER_IP> <DOMAIN_NAME>"
    exit 1
fi
echo "entrou porra2"

# Configurações
SERVER_USER="root"
SERVER_IP="$1"  # Usa o primeiro argumento como o IP do servidor
DOMAIN_NAME="$2"
SERVER_DIR="/root"
LOCAL_DIR="."  # Diretório do projeto
echo "entrou porra3"


# ssh-keygen -R "$SERVER_IP"
echo "serverIp: "$SERVER_IP
echo "domainName" $DOMAIN_NAME
echo "dirName" $SERVER_DIR
echo "Preparando diretório no servidor..."
echo $SERVER_USER@$SERVER_IP
# ssh -o StrictHostKeyChecking=no $SERVER_USER@$SERVER_IP

echo "Copiando arquivo global_portainer.yaml para o servidor..."
ssh $SERVER_USER@$SERVER_IP
# Substitui os placeholders nos arquivos .yaml
# sed "s/{{DOMAIN_NAME}}/$DOMAIN_NAME/g" global_portainer.yaml.template > global_portainer.yaml
echo "deu certo"
scp global_portainer.yaml $SERVER_USER@$SERVER_IP:$SERVER_DIR

# echo "Copiando arquivo traefik.yaml para o servidor..."
# scp global_traefik.yaml $SERVER_USER@$SERVER_IP:$SERVER_DIR

echo "Copiando arquivo rep_mongo.yaml para o servidor..."
scp rep_mongo.yaml $SERVER_USER@$SERVER_IP:$SERVER_DIR

echo "Copiando arquivo rep_evolution_api.yaml para o servidor..."
scp rep_evolution_api.yaml $SERVER_USER@$SERVER_IP:$SERVER_DIR

echo "Processando no servidor..."
ssh $SERVER_USER@$SERVER_IP <<EOF
    # Atualiza os pacotes
    echo "Atualizando pacotes do sistema..."
    sudo apt-get update -y

    # Instala o Docker se ele não estiver instalado
    if ! [ -x "\$(command -v docker)" ]; then
        echo "Instalando Docker..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
    fi

    # Iniciando o Swarm
    echo "Iniciando o Swarm"
    docker swarm init --advertise-addr=$SERVER_IP

    echo "criar rede para Traefik"
    docker network create --driver=overlay evolution_network

    echo "Executando o global_portainer.yaml"
    cd $SERVER_DIR

    docker stack deploy -c global_portainer.yaml portainer

    # docker stack deploy -c global_traefik.yaml traefik

    docker stack deploy -c rep_mongo.yaml postgres

    docker stack deploy -c rep_evolution_api.yaml evolution

EOF
