package integration

// example response json types

// provide a users current authentication status
type AuthStatus struct {
	// wether or not the user is authenticated
	IsAuthenticated bool `json:"is_authenticated"`
}

type CallbackStatus struct {
	Ok bool `json:"ok"`
}
