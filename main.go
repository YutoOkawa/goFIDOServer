package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

type helloJSON struct {
	UserName string `json:"username"`
}

type Config struct {
	ChallengeSize           int
	RpName                  string
	RpId                    string
	RpOrigin                string
	Timeout                 int
	RequireResidentKey      bool
	AuthenticatorAttachment string
	UserVerification        string
	Attestation             string
	CryptoParams            struct {
		Type string
		Alg  int
	}
}

var config Config

func init() {
	config.ChallengeSize = 64
	config.RpName = "FIDO_SERVER"
	config.RpId = "http://localhost:8080"
	config.RpOrigin = "http://localhost:8080"
	config.Timeout = 60000
	config.RequireResidentKey = false
	config.AuthenticatorAttachment = "cross-platform"
	config.UserVerification = "preferred"
	config.Attestation = "direct"
	config.CryptoParams.Type = "public-key"
	config.CryptoParams.Alg = -7
}

func main() {
	r := chi.NewRouter()
	r.Route("/attestation", func(r chi.Router) {
		r.Post("/options", attestationOptions)
	})
	http.ListenAndServe(":8080", r)
}
