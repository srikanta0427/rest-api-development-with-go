package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/srikanta0427/rest_api_design/model"
	"github.com/srikanta0427/rest_api_design/services"
	"github.com/srikanta0427/rest_api_design/validation"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// handler function

func (u *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	us := user.UserCons(6, "srikanta", "srikanta", "password", time.Now())

	// validate userStruct field
	errValid := validation.ValidateUserStruct(us)
	if errValid != nil {
		http.Error(w, errValid.Error(), http.StatusBadRequest)
		return
	}
	// UserService
	_, err := u.userService.CreateUser(us)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
}
