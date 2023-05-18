#!/bin/bash

set -e

echo "Remove previous stack"

# scaling down previous stack
docker service scale globalwebindex_blacklist=0
docker service scale globalwebindex_user=0
docker service scale globalwebindex_mysql=0
docker service scale globalwebindex_redis=0
docker service scale globalwebindex_vault=0

docker stack rm globalwebindex

# todo: wait until cleanup
sleep 20s

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