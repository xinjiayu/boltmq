package main

import (
	"git.oschina.net/cloudzone/smartgo/stgremoting/netm"
)

func main() {
	b := netm.NewBootstrap()
	b.Connect("10.122.1.200", 8000)
}
