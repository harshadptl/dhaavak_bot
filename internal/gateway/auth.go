package gateway

import (
	"crypto/subtle"
	"net/http"
)

// Authenticator validates connection tokens.
type Authenticator struct {
	token string
}

// NewAuthenticator creates a token authenticator. If token is empty, all connections are allowed.
func NewAuthenticator(token string) *Authenticator {
	return &Authenticator{token: token}
}

// Check validates the token from the request. Returns true if auth passes.
func (a *Authenticator) Check(r *http.Request) bool {
	if a.token == "" {
		return true
	}
	tok := r.URL.Query().Get("token")
	if tok == "" {
		tok = r.Header.Get("Authorization")
		if len(tok) > 7 && tok[:7] == "Bearer " {
			tok = tok[7:]
		}
	}
	return subtle.ConstantTimeCompare([]byte(tok), []byte(a.token)) == 1
}
