# Start Up Hyperledger Fabric Network

To start up an example Hyperledger Fabric network:

```bash
$ ./image-pull.sh
$ ./download-binaries.sh
$ export PATH=<path to current location>/bin:$PATH
$ cd start-network
$ ./byfn.sh -m generate
$ ./byfn.sh -m up
```

