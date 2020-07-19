package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Route("/attestation", func(r chi.Router) {
		r.Post("/options", attestationOptions)
	})
	// http.ListenAndServe(":8080", r)
	err := http.ListenAndServeTLS(":8080", "ssl/myself.crt", "ssl/myself.key", r)
	if err != nil {
		log.Fatal(err)
	}
}
