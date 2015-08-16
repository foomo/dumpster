package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/foomo/dumpster/config"
	"github.com/foomo/dumpster/server"
)

var flagConfig = flag.String("config", "", "where is my config")
var flagVersion = flag.Bool("version", false, "show version")

func main() {
	flag.Parse()
	if *flagVersion {
		fmt.Println("0.1.0")
		os.Exit(0)
	}
	if *flagConfig == "" {
		fmt.Println("i need a config -config /path/to/config.yml")
		//flag.PrintDefaults()
		os.Exit(1)
	}
	c, err := config.LoadConfig(*flagConfig)
	if err != nil {
		log.Fatal("could not start: " + err.Error())
	}
	err = server.Run(c)
	if err != nil {
		log.Fatal("exiting: " + err.Error())
	}
}
