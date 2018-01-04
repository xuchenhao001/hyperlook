package main

import (
	"github.com/xuchenhao001/Hyperlook/cadvisor"
	"time"
	"flag"
)

var argIp = flag.String("listen_ip", "", "IP to listen on, defaults to all IPs")
var argPort = flag.Int("port", 1001, "port to listen")

func main() {
	flag.Parse()

	cadvisor.New(argIp, argPort)

	time.Sleep(30*time.Minute)
}
