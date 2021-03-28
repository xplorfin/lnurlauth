package storage

import (
	"encoding/base64"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	lnurlHelper "github.com/fiatjaf/go-lnurl"
	. "github.com/stretchr/testify/assert"
	"github.com/xplorfin/lnurlauth"
)

func TestMemorySessionStore(t *testing.T) {
	newStore := MemorySessionStore{}

	key := gofakeit.Word()
	Nil(t, newStore.GetK1(key))

	sessionData := lnurlauth.SessionData{
		LnUrl: base64.StdEncoding.EncodeToString([]byte(gofakeit.Sentence(10))),
		Key:   gofakeit.BitcoinPrivateKey(),
	}

	randomKey := lnurlHelper.RandomK1()
	newStore.SetK1(randomKey, sessionData)

	res := newStore.GetK1(randomKey)
	Equal(t, res, &sessionData)

	newStore.RemoveK1(randomKey)
	Nil(t, newStore.GetK1(key))

	newStore.SetJwt(key, sessionData)
	res = newStore.GetJwt(key)
	Equal(t, res, &sessionData)

	newStore.RemoveJwt(key)
}
