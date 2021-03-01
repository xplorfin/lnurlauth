package main

import (
	"encoding/json"
	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/lnurlauth/integration"
	"github.com/xplorfin/netutils/testutils"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"
)

// run test with a timeout
func testWithTimeOut(t *testing.T, timeout <-chan time.Time, testFunc func(t *testing.T)) {
	done := make(chan bool)

	go func() {

		// do your testing here
		testFunc(t)

		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("test didn'tester finish in time")
	case <-done:
	}
}

func RunLnUrlTests(t *testing.T, url string) {
	testutils.AssertConnected(url, t)
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	Nil(t, err)
	client := &http.Client{
		Jar: jar,
	}

	runner := integration.TestRunner{
		Client: client,
		Tester: t,
		Url:    url,
	}

	False(t, runner.GetAuthStatus().IsAuthenticated)
	// get the url
	lnUrl, didRedirect := runner.GetLoginPage()
	// make sure we're still logged out
	False(t, runner.GetAuthStatus().IsAuthenticated)
	False(t, didRedirect)

	signedUrl, err := integration.SignCallbackUrl(lnUrl)
	Nil(t, err)

	// get without any context/cookie jar
	resp, err := http.Get(signedUrl)
	Nil(t, err)
	status := integration.CallbackStatus{}
	err = json.NewDecoder(resp.Body).Decode(&status)
	True(t, status.Ok)

	True(t, runner.GetAuthStatus().IsAuthenticated)
}

// TODO come up with a cleaner way to do this
func TestLocalTunnelsStart(t *testing.T) {
	testWithTimeOut(t, time.After(time.Second*20), func(t *testing.T) {
		args := os.Args[0:1]
		args = append(args, "-open=false")

		go Cmd(args)
		for {
			// wait till server url is set
			if serverUrl != "" {
				RunLnUrlTests(t, serverUrl)
				break
			}
			time.Sleep(time.Millisecond * 50)
		}

	})
}
