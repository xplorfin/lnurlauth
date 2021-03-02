package lnurlauth

type SessionData struct {
	LnUrl  string `json:"ln_url"`
	Key    string `json:"key"`
}

// persistent store for authenticating users cross-request
type SessionStore interface {
	// store session data
	Set(k1 string, value SessionData)
	// get data
	Get(k1 string) *SessionData
	// remove from session storage
	Remove(name string) *SessionData
}

// store for reading/storing data into the client
type RequestStore interface {
	// store request data
	Set(name, value string)
	// get data
	Get(name string) string
	// remove from request store, return removed value
	Remove(name string) string
}
