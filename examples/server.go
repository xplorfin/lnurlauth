package main

import (
	"context"
	"fmt"
	"github.com/localtunnel/go-localtunnel"
	"github.com/pkg/browser"
	"github.com/xplorfin/lnurlauth/integration"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
)

var serverUrl = ""

func Start(ctx context.Context, localTunnels, open bool, port, url string) error {
	var (
		server              http.Server
		localTunnelListener *localtunnel.Listener
		err                 error
	)

	// Setup a localTunnelListener for localtunnel
	if localTunnels {
		localTunnelListener, err = localtunnel.Listen(localtunnel.Options{})
		if err != nil {
			panic(err)
		}
		serverUrl = localTunnelListener.URL()
	} else {
		serverUrl = url
		if port != ""{
			serverUrl = fmt.Sprintf("%s:%s", serverUrl, port)
		}
	}

	server = integration.GenerateServer(serverUrl)

	g, _ := errgroup.WithContext(ctx)

	// Handle request from localtunnel
	g.Go(func() error {
		if localTunnels {
			fmt.Println(fmt.Sprintf("starting server at %s", serverUrl))
			err = server.Serve(localTunnelListener)
		} else {
			fmt.Println(fmt.Sprintf("starting server at %s on port %s", serverUrl, port))
			server.Addr = fmt.Sprintf(":%s", port)
			err = server.ListenAndServe()
		}
		if err != nil {
			panic(err)
		}
		return nil
	})
	g.Go(func() error {
		fmt.Println(fmt.Sprintf("server listening on %s", url))
		// bypass localtunnel authorization screen for this ip
		if open {
			fmt.Println("attempting to open browser")
			_ = browser.OpenURL(url)
		}
		return nil
	})

	err = g.Wait()
	return err
}

func getEnv(configVar, defaultVar string) (result string) {
	result = os.Getenv(configVar)
	if result == "" {
		return defaultVar
	}
	return result
}
