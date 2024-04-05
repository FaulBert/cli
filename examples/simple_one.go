package main

import (
	"fmt"
	"os"

	"github.com/nazhard/cli"
)

func main() {
	// initiate dummy app
	appFlag := []cli.Flag{
		&cli.StringFlag{
			Name:  "m",
			Value: "moe",
			Action: func(cli cli.Context) {
				fmt.Println("flag m value ", cli.String().Get("m"))
			},
		},
	}
	app := cli.App{
		Flags: appFlag,
		Action: func(ctx cli.Context) {
			fmt.Println("flag : ", ctx.String().Get("m"))
		},
	}
	app.Name = "uwe"
	cmdAlias := []string{"r", "rnu", "nur"}
	cmdFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "m",
			Value: "moe",
		},
	}

	cmd := &cli.Command{
		Name:        "run",
		Usage:       "",
		Alias:       cmdAlias,
		Flags:       cmdFlags,
		Description: "simply run",
		Action: func(ctx cli.Context) {
			flagValue := ctx.String().Get("m")
			fmt.Printf("run command invoked with m flag value %s \n", flagValue)
		},
	}

	// add command to app
	app.AddCommand(cmd)

	// if need more, just create new cli.Command
	// then app.AddCommand it

	fmt.Println(app.Run(os.Args))
}
