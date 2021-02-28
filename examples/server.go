package main

import (
	"fmt"
	"github.com/localtunnel/go-localtunnel"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"github.com/xplorfin/lnurlauth/integration"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
)

var serverUrl = ""

func Start(args []string) {
	app := &cli.App{
		Name:  "server",
		Usage: "run an lnurl-auth server",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "open",
				Usage: "wether or not to open the browser",
				Value: true,
			},
		},
		Action: func(c *cli.Context) error {
			// Setup a localTunnelListener for localtunnel
			localTunnelListener, err := localtunnel.Listen(localtunnel.Options{})
			if err != nil {
				panic(err)
			}

			g, _ := errgroup.WithContext(c.Context)

			server := integration.GenerateServer(localTunnelListener.URL())
			// Handle request from localtunnel
			g.Go(func() error {
				fmt.Println("starting server")
				err = server.Serve(localTunnelListener)
				if err != nil {
					panic(err)
				}
				return nil
			})
			g.Go(func() error {
				fmt.Println(fmt.Sprintf("server listening on %s", localTunnelListener.URL()))
				serverUrl = localTunnelListener.URL()
				// bypass localtunnel authorization screen for this ip
				fmt.Println("attempting to open browser")
				if c.Bool("open") {
					_ = browser.OpenURL(localTunnelListener.URL())
				}
				return nil
			})

			err = g.Wait()
			return err
		},
	}
	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Start(os.Args)
}
