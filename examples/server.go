package main

import (
	"context"
	"fmt"
	"github.com/localtunnel/go-localtunnel"
	"github.com/pkg/browser"
	"github.com/xplorfin/lnurlauth/integration"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func Start(ctx context.Context, localTunnels, open bool) error {
	var (
		server http.Server
		localTunnelListener *localtunnel.Listener
		url string
	)

	// Setup a localTunnelListener for localtunnel
	if localTunnels {
		localTunnelListener, err := localtunnel.Listen(localtunnel.Options{})
		if err != nil {
			panic(err)
		}
		url = localTunnelListener.URL()
		server = integration.GenerateServer(url)
	} else {

	}

	g, _ := errgroup.WithContext(ctx)

	// Handle request from localtunnel
	g.Go(func() error {
		fmt.Println("starting server")
		err := server.Serve(localTunnelListener)
		if err != nil {
			panic(err)
		}
		return nil
	})
	g.Go(func() error {
		fmt.Println(fmt.Sprintf("server listening on %s", localTunnelListener.URL()))
		// bypass localtunnel authorization screen for this ip
		fmt.Println("attempting to open browser")
		if open {
			_ = browser.OpenURL(url)
		}
		return nil
	})

	err := g.Wait()
	return err
}