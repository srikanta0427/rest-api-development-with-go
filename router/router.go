package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type routeHandler func(chi.Router)

func SetUpRouter(userRouter *UserRouter, organizerRouter *OrganizerRouter) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Authorization", "Content-Type", "X-CSRF-Token, Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Route("/", func(r chi.Router) {

		// for user
		r.Route("/user", func(r chi.Router) {
			userRouter.Register(r)
			userRouter.Profile(r)
		})

		// for organizer

		r.Route("/organizer", func(r chi.Router) {
			organizerRouter.Register(r)
			organizerRouter.Profile(r)
		})
	})

	return router

}
