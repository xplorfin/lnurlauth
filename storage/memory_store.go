package storage

import (
	"strings"
	"sync"

	"github.com/xplorfin/lnurlauth"
)

type MemorySessionStore struct {
	Map sync.Map
}

func (m *MemorySessionStore) get(key string) *lnurlauth.SessionData {
	loadedData, ok := m.Map.Load(key)
	if !ok {
		return nil
	}

	return loadedData.(*lnurlauth.SessionData)
}

func (m *MemorySessionStore) set(key string, value lnurlauth.SessionData) {
	m.Map.Store(key, &value)
}

func (m *MemorySessionStore) GetK1(k1 string) *lnurlauth.SessionData {
	key := checkAddPrefix(k1, lnurlauth.K1Prefix)

	return m.get(key)
}

func (m *MemorySessionStore) SetK1(k1 string, value lnurlauth.SessionData) {
	key := checkAddPrefix(k1, lnurlauth.K1Prefix)

	m.set(key, value)
}

func (m *MemorySessionStore) RemoveK1(k1 string) {
	key := checkAddPrefix(k1, lnurlauth.K1Prefix)

	m.Map.Delete(key)
}

func (m *MemorySessionStore) GetJwt(jwt string) *lnurlauth.SessionData {
	key := checkAddPrefix(jwt, lnurlauth.JwtPrefix)

	return m.get(key)
}

func (m *MemorySessionStore) SetJwt(jwt string, value lnurlauth.SessionData) {
	key := checkAddPrefix(jwt, lnurlauth.JwtPrefix)

	m.set(key, value)
}

func (m *MemorySessionStore) RemoveJwt(jwt string) {
	key := checkAddPrefix(jwt, lnurlauth.JwtPrefix)

	m.Map.Delete(key)
}

func (m *MemorySessionStore) Remove(key string) {
	m.Map.Delete(key)
}

func checkAddPrefix(key string, prefix string) string {
	if !strings.HasPrefix(key, prefix) {
		key = prefix + key
	}

	return key
}

var _ lnurlauth.SessionStore = &MemorySessionStore{}
