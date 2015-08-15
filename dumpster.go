package main

import (
	"github.com/foomo/dumpster/config"
	"github.com/foomo/dumpster/mock"
	"github.com/foomo/dumpster/server"
)

func main() {
	c, err := config.LoadConfig(mock.GetFilename("config.yml"))
	if err != nil {
		panic(err)
	}
	panic(server.Run(c))
}
