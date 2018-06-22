package main

import (
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/n3integration/conseil/actions"
)

var version = "dev"

func main() {
	log.SetFlags(0)
	log.SetPrefix("[conseil] ")

	app := cli.NewApp()
	app.Name = "conseil"
	app.Version = version
	app.Commands = actions.GetCommands()
	app.Usage = "go rapid application development"

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
