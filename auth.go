package lnurlauth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fiatjaf/go-lnurl"
	qrcode "github.com/skip2/go-qrcode"
)

// Authenticate authenticates a request,
// returns the user's LNUrl auth key if successful,
// returns an error otherwise.
func Authenticate(r *http.Request) (string, string, error) {
	var (
		k1, sig, key string
	)

	k1 = r.URL.Query().Get("k1")
	if k1 == "" {
		return "", "", errors.New("missing required key: k1")
	}

	sig = r.URL.Query().Get("sig")
	if sig == "" {
		return "", "", errors.New("missing required key: sig")
	}

	key = r.URL.Query().Get("key")
	if key == "" {
		return "", "", errors.New("missing required key: sig")
	}

	ok, err := lnurl.VerifySignature(k1, sig, key)
	if err != nil {
		return "", "", err
	} else if !ok {
		return "", "", errors.New("invalid signature")
	}

	return key, k1, nil
}

// GenerateQrCode generates a qr code from a full LNUrl.
func GenerateQrCode(lnurl string) ([]byte, error) {
	return qrcode.Encode(fmt.Sprintf("lightning:%s", lnurl), qrcode.Highest, 256)
}

// GenerateLnUrl generates a LNUrl string using a callback url
func GenerateLnUrl(callbackUrl string) (string, string, error) {
	generatedUrl := fmt.Sprintf("%s?k1=%s&tag=login", callbackUrl, lnurl.RandomK1())

	encoded, err := lnurl.LNURLEncode(generatedUrl)
	if err != nil {
		return encoded, "", err
	}

	return strings.ToLower(encoded), generatedUrl, nil
}
