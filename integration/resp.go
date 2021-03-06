package integration

// example response json types

// AuthStatus returns the current authentication status of the user.
type AuthStatus struct {
	// Whether or not the user is authenticated
	IsAuthenticated bool `json:"is_authenticated"`
}

// CallbackStatus returns whether or not a callback was successful to a ln wallet
type CallbackStatus struct {
	Ok bool `json:"ok"`
}
