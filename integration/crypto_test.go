package integration

import (
	"testing"

	"github.com/fiatjaf/go-lnurl"
	. "github.com/stretchr/testify/assert"
)

// verify integration signing works correctly
func TestSigning(t *testing.T) {
	var signingKey = ""
	for i := 0; i < 3; i++ {
		k1 := lnurl.RandomK1()
		res, err := signMessage(k1)
		Nil(t, err)
		ok, err := lnurl.VerifySignature(k1, res.Sig, res.Key)
		Equal(t, ok, true)
		Equal(t, res.Key, "02a673638cb9587cb68ea08dbef685c6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5")
		// make sure signing key is dynamic
		NotEqual(t, res.Sig, signingKey)
		signingKey = res.Sig
		Nil(t, err)
	}
}
