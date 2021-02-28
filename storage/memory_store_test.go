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
	Nil(t, newStore.Get(key))

	sessionData := lnurlauth.SessionData{
		LnUrl:  base64.StdEncoding.EncodeToString([]byte(gofakeit.Sentence(10))),
		RawUrl: gofakeit.URL(),
		Key:    gofakeit.BitcoinPrivateKey(),
	}

	randomKey := lnurlHelper.RandomK1()
	newStore.Set(randomKey, sessionData)

	res := newStore.Get(randomKey)
	Equal(t, res, &sessionData)

	newStore.Remove(randomKey)
	Nil(t, newStore.Get(key))
}
