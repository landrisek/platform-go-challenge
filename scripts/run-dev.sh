#!/bin/bash

set -e

echo "Build server"
sh scripts/build-server.sh

echo "Build local docker stack"
docker build -t server -f ./artifacts/server/Dockerfile .
#docker build -t server -f ./artifacts/migrations/Dockerfile .

set -a
source ./artifacts/.env
set +a
docker stack deploy -c ./artifacts/docker-compose.yml globalwebindex

sh scripts/wait-for-mysql.sh

echo "Provision vault"
# Generate random values for CRYPTO_KEY and HASHER_SALT
CRYPTO_KEY=$(openssl rand -hex 16)
HASHER_SALT=$(openssl rand -hex 8)

# Export the values as environment variables
export CRYPTO_KEY
export HASHER_SALT
sh scripts/provision-vault.sh