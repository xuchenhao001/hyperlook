#!/bin/bash

# clean k8s
kubectl delete deployment --all -n fabric-net
kubectl delete pods --all -n fabric-net
kubectl delete service --all -n fabric-net

# clean all your hyperledger/fabric running containers
docker ps -a | grep "dev" | awk '{ print $1 }' | xargs docker rm -fv
# clean all your chaincode images
docker images | grep "dev-peer" | awk '{ print $1 }' | xargs docker rmi -f

# clean fabric-net data
rm -rf /var/fabric-net/

