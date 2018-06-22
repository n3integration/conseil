package main

import (
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/n3integration/goji/actions"
)

var version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "goji"
	app.Usage = "go rapid application development"
	app.Commands = actions.Commands()
	app.Version = version

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
