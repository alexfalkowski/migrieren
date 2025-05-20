package token

import "github.com/alexfalkowski/go-service/v2/token"

// NewVerifier for token.
func NewVerifier(token *token.Token) token.Verifier {
	return token
}
