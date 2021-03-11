package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/localtunnel/go-localtunnel"
)

// generate a new local tunnel listener, store to file
func newLocalTunnelListener(f *os.File) (listener *localtunnel.Listener, err error) {
	listener, err = localtunnel.Listen(localtunnel.Options{})
	if err != nil {
		return nil, err
	}
	rawUrl := listener.URL()
	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Println("we failed to store listener url in cache")
		return listener, nil
	}
	subdomain := strings.Replace(u.Hostname(), ".loca.lt", "", 1)
	if f != nil {
		_, err = f.WriteString(subdomain)
		return listener, err
	}
	return listener, err
}

// since local tunnels has a phishing warning on every new domain (see: https://bit.ly/3kSn6Yp)
// we want to use the same domain from previous startups if we can
func getLocaltunnelsListener() (listener *localtunnel.Listener, err error) {
	path, err := os.Getwd()
	if err != nil {
		return newLocalTunnelListener(nil)
	}

	localTunnelsFileName := fmt.Sprintf("%s/.localtunnels", path)

	// try to open local tunnels file, create it if it doesn't exist
	f, err := os.OpenFile(localTunnelsFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		// if we can't find it just create a listener without saving it
		return newLocalTunnelListener(nil)
	}

	// read the contents of the local tunnels file for a subdomian
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return newLocalTunnelListener(f)
	}

	// if the contents don't exist, use create a new subdomani
	if len(contents) == 0 {
		return newLocalTunnelListener(f)
	} else {
		// otherwise, see if we can use the contents as a subdomain
		listener, err = localtunnel.Listen(localtunnel.Options{
			Subdomain:      string(contents),
			BaseURL:        "",
			MaxConnections: 0,
			Log:            nil,
		})
		if err != nil {
			log.Println(err)
			// if we can't, create a new listener and attempt to persist for next time
			return newLocalTunnelListener(f)
		}
	}
	return listener, err
}

func getEnv(configVar, defaultVar string) (result string) {
	result = os.Getenv(configVar)
	if result == "" {
		return defaultVar
	}
	return result
}
