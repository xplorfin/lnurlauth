package storage

import (
	"net/http"

	"github.com/xplorfin/lnurlauth"
)

// CookieRequestStore implements the request store interface and
// stores data in cookies using net/http
type CookieRequestStore struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func CookieStore(w http.ResponseWriter, r *http.Request) *CookieRequestStore {
	return &CookieRequestStore{
		Writer:  w,
		Request: r,
	}
}

func (c *CookieRequestStore) Set(name, value string) {
	http.SetCookie(c.Writer, &http.Cookie{Name: name, Value: value, HttpOnly: false})
}

func (c CookieRequestStore) Get(name string) string {
	storedData, err := c.Request.Cookie(name)
	if err != nil || storedData == nil {
		return ""
	}
	return storedData.Value
}

func (c *CookieRequestStore) Remove(name string) {
	http.SetCookie(c.Writer, &http.Cookie{Name: name, MaxAge: -1, Value: ""})
}

var _ lnurlauth.RequestStore = &CookieRequestStore{}
