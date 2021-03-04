package lnurlauth

const (
	// K1Prefix should be applied to all
	// key names when being stored or retrieved using the
	// SessionStore.GetK1() or SessionStore.SetK1() methods.
	K1Prefix string = "k1:"
	// JwtPrefix should be applied to all
	// key names when being stored or retrieved using the
	// SessionStore.GetJwt() or SessionStore.SetJwt() methods.
	JwtPrefix string = "jwt:"
)
