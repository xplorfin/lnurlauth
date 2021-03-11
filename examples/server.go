package main

import (
	"context"
	"fmt"

	"github.com/localtunnel/go-localtunnel"
	"github.com/pkg/browser"
	"github.com/xplorfin/lnurlauth/integration"
	"golang.org/x/sync/errgroup"
)

var serverUrl = ""

// start the server on a given port. If the localTunnels option is passed port and url are
// ignored in favor of a local tunnels url produced at runtime
func Start(ctx context.Context, localTunnels, open bool, port, url string) (err error) {
	// local tunnels listener. Set if local tunnels is in use
	var localTunnelListener *localtunnel.Listener

	if localTunnels {
		// attempt to create a local tunnels listener using cache
		localTunnelListener, err = getLocaltunnelsListener()
		if err != nil {
			panic(err)
		}
		serverUrl = localTunnelListener.URL()
	} else {
		// setup server on given host and port
		serverUrl = url
		if port != "" {
			url = fmt.Sprintf(":%s", port)
			serverUrl = fmt.Sprintf("http://%s:%s", "localhost", url)
		}
	}

	server := integration.GenerateServer()

	g, _ := errgroup.WithContext(ctx)

	// Handle request from localtunnel
	g.Go(func() error {
		if localTunnels {
			fmt.Printf("starting server at %s", serverUrl)

			err = server.Serve(localTunnelListener)
		} else {
			fmt.Printf("starting server at %s on port %s", serverUrl, port)

			server.Addr = url

			err = server.ListenAndServe()
		}
		if err != nil {
			panic(err)
		}
		return nil
	})
	g.Go(func() error {
		fmt.Printf("server listening on %s", serverUrl)
		// bypass localtunnel authorization screen for this ip
		if open {
			fmt.Println("attempting to open browser")
			_ = browser.OpenURL(serverUrl)
		}
		return nil
	})

	err = g.Wait()
	return err
}
