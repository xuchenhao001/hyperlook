package main

import (
	"fmt"
	"github.com/xuchenhao001/Hyperlook/cadvisor"
	"time"
	"flag"
)

var argIp = flag.String("listen_ip", "", "IP to listen on, defaults to all IPs")
var argPort = flag.Int("port", 8080, "port to listen")

func main() {
	fmt.Print("Hello, world!")

	cadvisor.New(argIp, argPort)

	time.Sleep(30*time.Minute)
}
