package controller

import (
	"encoding/json"
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

	// json -> go struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user = user.UserCons(9, user.Name, user.Email, user.Password, time.Now())

	// validate userStruct field
	errValid := validation.ValidateUserStruct(user)
	if errValid != nil {
		http.Error(w, errValid.Error(), http.StatusBadRequest)
		return
	}

	// UserService
	_, errSer := u.userService.CreateUser(&user)
	if errSer != nil {
		fmt.Println(errSer.Error())
		panic(errSer.Error())
	}
	w.WriteHeader(http.StatusCreated)
}
