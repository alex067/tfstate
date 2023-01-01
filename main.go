/*
Copyright © 2022 NAME HERE panayi067@gmail.com

*/
package main

import (
	"alex067/tfstate/cmd"
	"log"
)

var version string

func confLogger() {
	log.SetFlags(log.Ltime)
}

func main() {
	confLogger()
	cmd.Version = version
	cmd.Execute()
}
