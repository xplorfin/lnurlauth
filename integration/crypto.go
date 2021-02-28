package integration

import (
	"encoding/hex"
	"fmt"
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
		fmt.Println(err)
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

func SignCallbackUrl(rawCallbackUrl string) (res string, err error) {
	callbackUrl, err := url.Parse(rawCallbackUrl)
	if err != nil {
		return res, err
	}
	challenge := callbackUrl.Query().Get("k1")
	signedMessage, err := signMessage(challenge)
	if err != nil {
		return res, err
	}
	// add query params
	q := callbackUrl.Query()
	q.Set("sig", signedMessage.Sig)
	q.Set("key", signedMessage.Key)
	// mutate callback url w/ new params
	callbackUrl.RawQuery = q.Encode()

	return callbackUrl.String(), nil
}
