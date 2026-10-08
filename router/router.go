package router

import (
	"github.com/go-chi/chi"
)

type routeHandler func(chi.Router)

func SetUpRouter(userRouter *UserRouter, organizerRouter *OrganizerRouter) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {

		// for user
		r.Route("/user", func(r chi.Router) {
			userRouter.Register(r)
		})

		// for organizer

		r.Route("/organizer", func(r chi.Router) {
			organizerRouter.Register(r)
		})
	})

	return router

}
