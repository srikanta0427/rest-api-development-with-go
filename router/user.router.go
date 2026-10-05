package router

import (
	"github.com/go-chi/chi"
	"github.com/srikanta0427/rest_api_design/controller"
)

type UserRouter struct {
	userController *controller.UserController
}

func NewUserRouter(userController *controller.UserController) Router {
	return &UserRouter{userController: userController}
}

func (u *UserRouter) Register(r chi.Router) {
	r.Post("/signup", u.userController.RegisterUser)
}
