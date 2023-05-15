#!/bin/bash

set -e

echo "Remove previous stack"
docker stack rm globalwebindex
# todo: wait until cleanup
sleep 20s

echo "Build server"
sh scripts/build-server.sh

set -a
source ./artifacts/.env
set +a

echo "Build local docker stack"
# docker rmi server
docker rmi migrations -f
docker build -t server -f ./artifacts/server/Dockerfile .
docker build -t migrations -f ./artifacts/migrations/Dockerfile .

docker stack deploy -c ./artifacts/docker-compose.yml globalwebindex