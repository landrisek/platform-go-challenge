#!/bin/bash

set -e

# scaling down previous stack
stack="globalwebindex"

if docker stack ls --format '{{.Name}}' | grep -q "^$stack$"; then
    echo "Removing previous '$stack' stack."
    docker service scale "$stack"_blacklist=0
    docker service scale "$stack"_user=0
    docker service scale "$stack"_mysql=0
    docker service scale "$stack"_redis=0
    docker service scale "$stack"_vault=0

    docker stack rm "$stack"
    # wait until cleanup
    sleep 20s
fi

set -a
source ./artifacts/.env
set +a

echo "Build local docker stack"
# docker rmi user
# docker rmi migrations -f
# docker rmi blacklist -f
docker build -t user -f ./artifacts/user/Dockerfile .
docker build -t migrations -f ./artifacts/migrations/Dockerfile .
docker build -t blacklist -f ./artifacts/blacklist/Dockerfile .

docker stack deploy -c ./artifacts/docker-compose.yml globalwebindex