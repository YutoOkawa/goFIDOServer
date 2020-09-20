package webauthn

import (
	"crypto/rand"
	"encoding/base64"
)

func makeRandom(i int) string {
	b := make([]byte, i)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
