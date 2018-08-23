# Start Up Hyperledger Fabric Network

To start up an example Hyperledger Fabric network:

```bash
cd hyperlook/fabric-sample-network/
./image-pull.sh
./download-binaries.sh
cp bin/* /usr/bin/
kubectl create ns fabric-net
kubectl label node <your-host> bc=true
./start-network.sh
```

Tear down your network:

```bash
./clean-network.sh
```
