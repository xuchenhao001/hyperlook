#!/bin/bash

cd network
# generate certs
./generateCerts.sh

# start new network
kubectl create -f manifests/
