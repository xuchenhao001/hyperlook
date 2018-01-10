package main

import (
	"github.com/xuchenhao001/Hyperlook/blockchainPerform/fabricSDK"
	"github.com/xuchenhao001/Hyperlook/cadvisor"
	"flag"
)

var argIp = flag.String("listen_ip", "", "IP to listen on, defaults to all IPs")
var argPort = flag.Int("port", 2053, "port to listen")

func main() {
	flag.Parse()

	go func() {
		cadvisor.New(argIp, argPort)
	}()

	fabricSDK.New()

}
