#!/usr/bin/env bash
MODULE="smpp_simulator"

./build.sh

echo ">> SMPP: start new container"
net_opt="--network=host"
if [ "$1" == "bridge" ]; then
    net_opt=""
fi

docker run -d -p 2775:2775 -p 2777:2777 -p 2779:2779 $net_opt --name $MODULE $MODULE
sleep 5

echo ">> SMPP: done"
docker ps -f name=$MODULE
