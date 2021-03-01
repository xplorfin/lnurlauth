package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var serverUrl = ""

func Cmd(args []string) {
	app := &cli.App{
		Name:  "server",
		Usage: "run an lnurl-auth server",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "open",
				Usage: "wether or not to open the browser",
				Value: true,
			},
			&cli.BoolFlag{
				Name:        "localtunnels",
				Usage:       "use local tunnels (useful for testing w/ a wallet)",
				Value:       true,
			},
		},
		Action: func(c *cli.Context) error {
			return Start(c.Context, c.Bool("localtunnels"), c.Bool("open"))
		},
	}
	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Cmd(os.Args)
}
