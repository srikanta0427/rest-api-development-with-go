package router

import (
	"github.com/go-chi/chi"
	"github.com/srikanta0427/rest_api_design/controller"
)

type Router interface {
	Register(r chi.Router)
}

func SetUpRouter(userRouter Router) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", controller.RunGet)
	router.Get("/{id}", controller.RunPost)

	userRouter.Register(router)

	return router
}
