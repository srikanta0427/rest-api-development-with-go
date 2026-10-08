package controller

import (
	"encoding/json"
	"net/http"

	"github.com/srikanta0427/rest_api_design/model"
	"github.com/srikanta0427/rest_api_design/services"
	"github.com/srikanta0427/rest_api_design/util"
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

	// while returning user to client
	user.Password = "Thank You😊😊😊"

	// validate userStruct field
	errValid := validation.ValidateUserStruct(user)
	if errValid != nil {
		util.Error(w, http.StatusBadRequest, "validation failed", errValid.Error())
		return
	}



	// UserService
	_, errSer := u.userService.CreateUser(user)
	if errSer != nil {
		util.Error(w, http.StatusBadRequest, "registration failed", errSer.Error())
	}
	util.Success(w, http.StatusCreated, "register success", user)
}
