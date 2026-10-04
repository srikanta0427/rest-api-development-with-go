package controller

import (
	"fmt"
	"net/http"

	"github.com/srikanta0427/rest_api_design/services"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	_ = u.userService.CreateUser()
	fmt.Println("CreateUser from controller")
	w.Write([]byte("CreateUser from controller"))
}
