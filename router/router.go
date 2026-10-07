package router

import (
	"github.com/go-chi/chi"
)

type Router interface {
	Register(r chi.Router)
}

func SetUpRouter(route Router) *chi.Mux {
	router := chi.NewRouter()

	route.Register(router)
	return router
}
