package storage

import (
	"sync"

	"github.com/xplorfin/lnurlauth"
)

type MemorySessionStore struct {
	Map sync.Map
}

func (m *MemorySessionStore) Set(k1 string, value lnurlauth.SessionData) {
	m.Map.Store(k1, &value)
}

func (m *MemorySessionStore) Get(k1 string) *lnurlauth.SessionData {
	loadedData, ok := m.Map.Load(k1)
	if !ok {
		return nil
	}
	return loadedData.(*lnurlauth.SessionData)
}

func (m *MemorySessionStore) Remove(name string) *lnurlauth.SessionData {
	toRemove := m.Get(name)
	m.Map.Delete(name)
	return toRemove
}

var _ lnurlauth.SessionStore = &MemorySessionStore{}
