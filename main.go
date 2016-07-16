package main

import (
	"os"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/raintank/snap-plugin-collector-tcpconns/tcpconns"
)

func main() {

	plugin.Start(
		tcpconns.Meta(),
		tcpconns.New(),
		os.Args[1],
	)
}
