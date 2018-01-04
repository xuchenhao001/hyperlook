package main

import (
	"fmt"
	"github.com/xuchenhao001/Hyperlook/cadvisor"
	"time"
)

func main() {
	fmt.Print("Hello, world!")
	var ip * string
	*ip = "0.0.0.0"
	cadvisor.New(ip, 1099)

	time.Sleep(30*time.Minute)
}
