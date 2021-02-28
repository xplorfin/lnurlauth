package storage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	. "github.com/stretchr/testify/assert"
)

func (c *CookieRequestStore) CopyHeaders() {
	c.Request = &http.Request{Header: http.Header{"Cookie": c.Writer.(*httptest.ResponseRecorder).HeaderMap["Set-Cookie"]}}
}

func TestCookieStore(t *testing.T) {
	request := &http.Request{}
	writer := &httptest.ResponseRecorder{Code: http.StatusOK}
	cookieStore := CookieRequestStore{
		writer,
		request,
	}

	key := gofakeit.Word()
	value := gofakeit.Word()

	// make sure we're not getting garbage data
	Equal(t, cookieStore.Get(key), "")
	cookieStore.Set(key, value)
	// Copy the Cookie over to a new Request
	cookieStore.CopyHeaders()
	// make sure we can get cookie data
	Equal(t, cookieStore.Get(key), value)
	cookieStore.Remove(key)
	// Copy the Cookie over to a new Request
	cookieStore.CopyHeaders()
	// since the cookies driver sets a max age header,
	// cookies won't be removed until they hit the browser
	// So we make sure the set-cookie
	Equal(t, len(cookieStore.Writer.(*httptest.ResponseRecorder).HeaderMap["Set-Cookie"]), 2)
}
