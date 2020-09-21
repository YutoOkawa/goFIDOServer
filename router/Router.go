package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Server struct {
	Router *chi.Mux
}

func New() *Server {
	return &Server{
		Router: chi.NewRouter(),
	}
}

func (s *Server) SetRouter() {
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.Router.Route("/attestation", func(register chi.Router) {
		register.Post("/options", AttestationOptions)
		register.Post("/result", AttestationResult)
	})

	s.Router.Route("/assertion", func(auth chi.Router) {
		auth.Post("/options", AssertionOptions)
	})
}
