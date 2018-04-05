package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Name = "solutions"
	app.Usage = "Solutions service"
	app.Flags = flags
	app.Action = initServer

	fmt.Printf("Starting %s %s\n", app.Name, app.Version)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
