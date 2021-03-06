package integration

import (
	"encoding/hex"
	"net/url"

	"github.com/btcsuite/btcd/btcec"
)

// a signed k1
type SignedMessage struct {
	// (hex encoded 32 bytes of challenge) which is going to be signed by user's linkingPrivKey.
	K1 string
	// signature of message using secp256k1 on the linkingPrivKey
	Sig string
	//  the linking key https://git.io/JqT7T
	Key string
}

// simulate a k1 signing
func signMessage(k1 string) (params SignedMessage, err error) {
	// sign the message with a pre-defined private key https://git.io/JqT7t
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c")

	if err != nil {
		return
	}

	// generate a public private key pair from the private-key
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	// decode the challenge from hex
	decodedChallenge, err := hex.DecodeString(k1)
	if err != nil {
		return params, err
	}

	// sign the challenge
	sig, err := privKey.Sign(decodedChallenge)
	if err != nil {
		return params, err
	}

	// return the signed message object
	return SignedMessage{
		K1:  k1,
		Sig: hex.EncodeToString(sig.Serialize()),
		Key: hex.EncodeToString(pubKey.SerializeCompressed()),
	}, nil
}

// sign a callback url with a pre-defined private key
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
