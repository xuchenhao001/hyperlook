# Start Up Hyperledger Fabric Network

To start up an example Hyperledger Fabric network:

```bash
$ ./image-pull.sh
$ ./download-binaries.sh
$ cp bin/* /usr/bin/
$ cd start-network
$ ./generateCerts.sh
$ kubectl create ns fabric-net
$ kubectl label node <fabric-net-node> bc=true
$ kubectl create -f manifest/
```

To operate your fabric network:

```bash
$ docker ps |grep netctl
$ docker exec -ti <your-netctl-containerID> sh
$ netctl createchannel 
$ netctl joinchannel 
$ netctl installCC 
$ netctl instantiateCC 
$ netctl invokeCC 
$ netctl queryCC 
$ netctl upgradeCC
```

Tear down your network:

```bash
$ ./clean-network.sh
```
