# Start Up Hyperledger Fabric Network

To start up an example Hyperledger Fabric network:

```bash
$ ./image-pull.sh
$ ./download-binaries.sh
$ export PATH=<path to current location>/bin:$PATH
$ cd start-network
$ ./generateCerts.sh
$ kubectl create ns fabric-net
$ kubectl label node <fabric-net-node> bc=true
$ kubectl create -f manifest/
```

