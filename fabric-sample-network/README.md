# Start Up Hyperledger Fabric Network

To start up an example Hyperledger Fabric network:

```bash
$ ./image-pull.sh
$ ./download-binaries.sh
$ cp bin/* /usr/bin/
$ kubectl create ns fabric-net
$ kubectl label node <fabric-net-node> bc=true
$ ./start-network.sh
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

# OR

$ docker ps |grep cli
$ docker exec -ti <your-cli-containerID> bash
$ scripts/script.sh
```

Tear down your network:

```bash
$ ./clean-network.sh
```
