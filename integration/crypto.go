package integration

import (
	"encoding/hex"
	"net/url"

	"github.com/btcsuite/btcd/btcec"
)

type SignedMessage struct {
	K1  string
	Sig string
	Key string
}

// simulate a k1 signing
func signMessage(k1 string) (params SignedMessage, err error) {
	// sign the message with a pre-defined private key
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c")

	if err != nil {
		return
	}

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	decodedChallenge, err := hex.DecodeString(k1)
	if err != nil {
		return params, err
	}

	sig, err := privKey.Sign(decodedChallenge)
	if err != nil {
		return params, err
	}

	return SignedMessage{
		K1:  k1,
		Sig: hex.EncodeToString(sig.Serialize()),
		Key: hex.EncodeToString(pubKey.SerializeCompressed()),
	}, nil
}

func SignCallbackUrl(rawCallbackUri string) (res string, err error) {
	callbackUri, err := url.Parse(rawCallbackUri)
	if err != nil {
		return res, err
	}

	challenge := callbackUri.Query().Get("k1")
	signedMessage, err := signMessage(challenge)
	if err != nil {
		return res, err
	}

	// add query params
	q := callbackUri.Query()
	q.Set("sig", signedMessage.Sig)
	q.Set("key", signedMessage.Key)

	// mutate callback url with new params
	callbackUri.RawQuery = q.Encode()

	return callbackUri.String(), nil
}
