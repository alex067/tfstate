/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"alex067/tfstate/cmd"
	"log"
)

func confLogger() {
	log.SetFlags(log.Ltime)
}

func main() {
	confLogger()
	cmd.Execute()
}
