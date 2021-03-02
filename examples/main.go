package main

import (
	"github.com/phayes/freeport"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
)


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
				Name:  "localtunnels",
				Usage: "use local tunnels (useful for testing w/ a wallet)",
				Value: true,
			},
			&cli.StringFlag{
				Name:  "port",
				Usage: "port to start the server on, ignored if localtunnels is set",
				Value: getEnv("PORT", strconv.Itoa(freeport.GetPort())),
			},
			&cli.StringFlag{
				Name:  "hostname",
				Usage: "hostname to start server on",
				Value: getEnv("rando", "localhost"),
			},
		},
		Action: func(c *cli.Context) error {
			return Start(c.Context, c.Bool("localtunnels"), c.Bool("open"), c.String("port"), c.String("url"))
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
