#!/bin/bash

kubectl delete ns fabric-net

# clear all your hyperledger/fabric running containers
docker ps -a | grep "chaincode\|fabric" | awk '{ print $1 }' | xargs docker rm -fv
# clear all your chaincode images
docker images | grep "dev-peer" | awk '{ print $1 }' | xargs docker rmi -f

# clear fabric-net data
rm -rf /var/fabric-net/

