package lnurlauth

type SessionData struct {
	LnUrl string `json:"ln_url"`
	Key   string `json:"key"`
}

// SessionStore is a persistent store for authenticating users cross-request.
// Note: for methods such as Get/Set K1/Jwt, it might be a good idea
// to prefix stored keys using something like "k1:[key]" or
// "jwt:[key]" to prevent possible collisions.
type SessionStore interface {
	// GetK1 should return lnurlauth
	// session data using a k1 challenge string
	// as the retrieval key.
	GetK1(string) *SessionData
	// SetK1 should store initial lnurlauth session data
	// using a k1 challenge string as its key.
	SetK1(string, SessionData)
	// RemoveK1 should remove lnurlauth
	// session data associated with a k1
	// challenge string key.
	RemoveK1(string)
	// GetJwt should retrieve lnurlauth
	// session data using a JSON web token
	// as the retrieval key.
	GetJwt(string) *SessionData
	// SetJwt should store lnurlauth
	// session data using a JSON web token as its key.
	SetJwt(string, SessionData)
	// RemoveJwt should remove lnurlauth
	// session data associated with a
	// JSON web token key.
	RemoveJwt(string)
	// Remove should remove session data using
	// any possible key (such as a K1 challenge or JSON web token).
	// Note: if your implementations of SetK1 or SetJwt apply prefixes
	// before their storage calls, the key you want to remove
	// should also be prefixed.
	Remove(string)
}

// RequestStore is a store for reading/storing data into the client
type RequestStore interface {
	// Set should take a key as its first param, and a
	// string (or stringified) value as its second param,
	// and store the data in whatever key-value store is being
	// used as a backend.
	Set(string, string)
	// Get should return
	Get(string) string
	// Remove should remove data from the request store using any possible key.
	Remove(string)
}
