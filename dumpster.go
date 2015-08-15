package main

import (
	"flag"

	"github.com/foomo/dumpster/config"
	"github.com/foomo/dumpster/server"
)

var flagConfig = flag.String("config", "", "where is my config")

func main() {
	flag.Parse()
	c, err := config.LoadConfig(*flagConfig)
	if err != nil {
		panic(err)
	}
	panic(server.Run(c))
}
