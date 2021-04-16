package main

import (
	"flag"

	"gitlab.com/jebo87/makako-gateway/app"
)

//Environment variables to set:
// "API_ADDRESS_PROD"
// "API_PORT_PROD")
// "API_ADDRESS_DEV"
// "API_PORT_DEV"
// "PORT"
// "ELASTIC_ADDRESS"
// "ELASTIC_PORT"

var deployedFlag *bool

func main() {
	loadConfig()
	app.StartApp()
}

func loadConfig() {
	//TODO: check if we really need to do this
	deployedFlag = flag.Bool("deployed", false, "Defines if absolute paths need to be used for the config files")
	flag.Parse()
}
