package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/srikanta0427/rest_api_design/controller"
)

type UserRouter struct {
	userController *controller.UserController
}

func NewUserRouter(userController *controller.UserController) *UserRouter {
	return &UserRouter{userController: userController}
}

func (u *UserRouter) Register(r chi.Router) {
	r.Post("/org-signup", u.userController.RegisterUser)
}

func (u *UserRouter) Profile(r chi.Router) {
	r.Get("/home", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("home"))
	})
}
