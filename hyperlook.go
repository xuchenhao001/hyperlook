package main

import (
	"fmt"
	"github.com/xuchenhao001/Hyperlook/cadvisor"
	"time"
)

func main() {
	fmt.Print("Hello, world!")
	var a cadvisor.ImageFsInfoProvider
	cadvisor.New("0.0.0.0", 1099, a,"/", false)

	time.Sleep(300*time.Second)
}
