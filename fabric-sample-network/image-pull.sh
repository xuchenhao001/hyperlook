#!/bin/bash

VERSION=1.0.5
ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
#Set MARCH variable i.e ppc64le,s390x,x86_64,i386
MARCH=`uname -m`
FABRIC_TAG="$MARCH-$VERSION"

dockerFabricPull() {
  local FABRIC_TAG=$1
  for IMAGES in peer orderer ca couchdb ccenv javaenv kafka zookeeper tools; do
      echo "==> FABRIC IMAGE: $IMAGES"
      echo
      docker pull hyperledger/fabric-$IMAGES:$FABRIC_TAG
      docker tag hyperledger/fabric-$IMAGES:$FABRIC_TAG hyperledger/fabric-$IMAGES
  done
}

echo "===> Pulling fabric Images"
dockerFabricPull ${FABRIC_TAG}
