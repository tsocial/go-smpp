#!/usr/bin/env bash
MODULE="smpp_simulator"

# clean up existing docker smpp sim server
if [ "$(docker ps -a -q -f name=$MODULE)" ]; then
    set +e
    
    echo ">> SMPP: stop old container"
    docker stop $MODULE

    echo ">> SMPP: remove old container"
    docker rm $MODULE

    set -e
fi

if [ "$(docker images -a -q $MODULE)" ]; then
    echo ">> SMPP: remove old image"
    docker rmi $MODULE:latest
fi

echo ">> SMPP: build new image"
docker build -f ./Dockerfile -t $MODULE .
