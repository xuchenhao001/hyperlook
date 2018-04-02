package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	peerJoinChannel = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "peer_join_channel_time_seconds",
		Help: "peer join channel time",
	})
	peerInstallChaincode = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "peer_install_chaincode_time_seconds",
		Help: "peer install chaincode time",
	})
	peerInstantiateChaincode = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "peer_instantiate_chaincode_time_seconds",
		Help: "peer instantiate chaincode time",
	})
	peerInvokeChaincode = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "peer_invoke_chaincode_time_seconds",
		Help: "peer invoke chaincode time",
	})
	peerQueryChaincode = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "peer_query_chaincode_time_seconds",
		Help: "peer query chaincode time",
	})
	peerUpgradeChaincode = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "peer_upgrade_chaincode_time_seconds",
		Help: "peer upgrade chaincode time",
	})
)

func init() {
	prometheus.MustRegister(peerJoinChannel)
	prometheus.MustRegister(peerInstallChaincode)
	prometheus.MustRegister(peerInstantiateChaincode)
	prometheus.MustRegister(peerInvokeChaincode)
	prometheus.MustRegister(peerQueryChaincode)
	prometheus.MustRegister(peerUpgradeChaincode)
}
