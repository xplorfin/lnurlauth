package lnurlauth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fiatjaf/go-lnurl"
	qrcode "github.com/skip2/go-qrcode"
)

// authenticate a request, return key if successful otherwise an error
func Authenticate(r *http.Request) (key, k1 string, err error) {
	k1 = r.URL.Query().Get("k1")
	if k1 == "" {
		return key, k1, errors.New("missing required key: k1")
	}
	sig := r.URL.Query().Get("sig")
	if sig == "" {
		return key, k1, errors.New("missing required key: sig")
	}
	key = r.URL.Query().Get("key")
	if key == "" {
		return key, k1, errors.New("missing required key: sig")
	}
	ok, err := lnurl.VerifySignature(k1, sig, key)
	if err != nil {
		return key, k1, err
	}
	if !ok {
		return key, k1, errors.New("invalid signature")
	}
	return key, k1, err
}

// generate a qr code from an lnurl
func GenerateQrCode(lnurl string) ([]byte, error) {
	png, err := qrcode.Encode(fmt.Sprintf("lightning:%s", lnurl), qrcode.Highest, 256)
	if err != nil {
		return nil, err
	}
	return png, nil
}

// generate an ln url string
func GenerateLnUrl(callbackUrl string) (encoded string, rawurl string, err error) {
	res := fmt.Sprintf("%s?k1=%s&tag=login", callbackUrl, lnurl.RandomK1())
	encoded, err = lnurl.LNURLEncode(res)
	if err != nil {
		return encoded, res, err
	}
	return strings.ToLower(encoded), res, err
}
