package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	peerJoinChannel = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "peer_join_channel_time_seconds",
		Help: "peer join channel time",
	},[]string{"container_name"})
	peerInstallChaincode = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "peer_install_chaincode_time_seconds",
		Help: "peer install chaincode time",
	},[]string{"container_name"})
	peerInstantiateChaincode = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "peer_instantiate_chaincode_time_seconds",
		Help: "peer instantiate chaincode time",
	},[]string{"container_name"})
	peerInvokeChaincode = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "peer_invoke_chaincode_time_seconds",
		Help: "peer invoke chaincode time",
	},[]string{"container_name"})
	peerQueryChaincode = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "peer_query_chaincode_time_seconds",
		Help: "peer query chaincode time",
	},[]string{"container_name"})
	peerUpgradeChaincode = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "peer_upgrade_chaincode_time_seconds",
		Help: "peer upgrade chaincode time",
	},[]string{"container_name"})
)

func init() {
	prometheus.MustRegister(peerJoinChannel)
	prometheus.MustRegister(peerInstallChaincode)
	prometheus.MustRegister(peerInstantiateChaincode)
	prometheus.MustRegister(peerInvokeChaincode)
	prometheus.MustRegister(peerQueryChaincode)
	prometheus.MustRegister(peerUpgradeChaincode)
}
