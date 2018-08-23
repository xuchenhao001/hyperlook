# Hyperlook

Hyperledger Fabric monitoring project based on kubernetes.

This project is meant to monitor your Fabric network's status. The monitoring targets have 2 part:

* Fabric peers' running status
* Fabric peers' transaction performance (join channel, install chaincode, instantiate chaincode, invoke, and query)

![](https://github.com/xuchenhao001/hyperlook/blob/master/images/peerRunningStatus.png)

![](https://github.com/xuchenhao001/hyperlook/blob/master/images/peerTransactionPerformance.png)

There are 5 parts of this project:

* hyperlook -> for peer's transaction performance monitoring
* netctl -> Fabric network controler
* grafana-dashboard -> monitoring GUI based on grafana
* prometheus-config -> alert confuration for prometheus
* fa8ric -> Fabric integration with Kubernetes (Just a workaround now)

## Prerequisite

* Ubuntu 16.04
* Memory more than 8GB

## Build

Build hyperlook binary and docker image:

```bash
cd hyperlook/
CGO_ENABLED=0 go build
docker build -t hyperlook:latest .
```

Build netctl binary and docker image:

```bash
cd hyperlook/fabric-sample-network/netctl/
./build.sh
```

## Deploy monitoring system

### 1. Start Kubernetes

I perfer to install Kubernetes through `kubeadm`, See [this official doc](https://kubernetes.io/docs/setup/independent/install-kubeadm/).

### 2. Deploy prometheus

I perfer to install prometheus through `prometheus-operator`, see [this doc](https://github.com/coreos/prometheus-operator/tree/master/contrib/kube-prometheus#quickstart).

### 3. Deploy Elasticsearch

Hyperlook can analysis logs based on elasticsearch, so we must deploy it first. See [this doc](https://github.com/kubernetes/kubernetes/blob/master/cluster/addons/fluentd-elasticsearch/README.md).

After deploy, the elasticsearch needs some time to work well (Maybe more than 20 mins, which is depends on your vm's performance).

### 4. Start up Hyperledger Fabric network

See [here](https://github.com/xuchenhao001/hyperlook/blob/master/fabric-sample-network/README.md).

### 5. Start up hyperlook

```bash
cd hyperlook/
kubectl create -f kubernetes/manifests/
```

### 6. Setup grafana dashboard

Open your browser and go to `http://<your-host-ip>:30902`, login with username `admin` and password `admin`. Then import my Fabric dashboard in `hyperlook/grafana-dashboard/Fabric.json`

### 7. Setup Alert rules

```bash
cd hyperlook/alert/
kubectl create -f prometheus-config.yaml
# make changes effective
curl -X POST <your-host-ip>:30900/-/reload
```

### 8. Operate your fabric network

```bash
docker exec -ti netctl sh
netctl createchannel 
netctl joinchannel 
netctl installCC 
netctl instantiateCC 
netctl invokeCC 
netctl queryCC 
netctl upgradeCC

# OR

docker exec -ti cli bash
scripts/script.sh
```

After operation, you will see the monitoring data curve on your grafana dashboard.

