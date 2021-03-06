package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	. "github.com/stretchr/testify/assert"
)

// helper methods
type TestRunner struct {
	// http client to use for tests
	Client *http.Client
	// testing object to report errors to
	Tester *testing.T
	// base url to test against
	Url string
}

// get a given route with a bypass header for local tunnels
func (t TestRunner) Get(route string) *http.Response {
	body := bytes.NewReader([]byte{})
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", t.Url, route), body)
	Nil(t.Tester, err)
	req.Header.Set("Bypass-Tunnel-Reminder", "don't show screen")
	res, err := t.Client.Do(req)
	Nil(t.Tester, err)
	return res
}

// retreive the users authentication status
func (t TestRunner) GetAuthStatus() AuthStatus {
	authStatus := AuthStatus{}
	res := t.Get("is-authenticated")
	err := json.NewDecoder(res.Body).Decode(&authStatus)
	Nil(t.Tester, err)
	return authStatus
}

// returns lnurl if availalable, and wether or not a recirect was attempted
func (t TestRunner) GetLoginPage() (lnurl string, didRedirect bool) {
	res := t.Get("login")
	parsedUrl, _ := url.Parse(t.Url)
	cookies := t.Client.Jar.Cookies(parsedUrl)
	for _, cookie := range cookies {
		if cookie.Name == CookieName {
			lnurl = cookie.Value
		}
	}
	if res.StatusCode == 301 || res.StatusCode == 302 {
		didRedirect = true
	}
	return lnurl, didRedirect
}
