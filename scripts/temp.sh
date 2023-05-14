#!/bin/bash

set -e

set -a
source ./artifacts/.env
set +a

docker rmi migrations -f
docker build -t migrations -f ./artifacts/migrations/Dockerfile .
docker run -e VAULT_TOKEN=myroot -e VAULT_ADDR=http://localhost:8200 -e VAULT_MOUNT="database" \
    -e MYSQL_DATABASE=sandbox -e MYSQL_USER=sandbox -e MYSQL_HOST=mysql -e MYSQL_PORT=3306 migrations 