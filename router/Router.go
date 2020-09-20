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
		// register.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// 	if err := db.InsertDB("aa", "test"); err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	if err := db.InsertDB("aaa", "test1"); err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	user, err := db.GetOneDB("aaa")
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	w.Write([]byte(user.UserID))
		// 	// if err := db.DeleteDB("aa"); err != nil {
		// 	// 	log.Fatal(err)
		// 	// }
		// })
	})
}
