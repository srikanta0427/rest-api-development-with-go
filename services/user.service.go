package services

import (
	"fmt"

	db "github.com/srikanta0427/rest_api_design/db/repos"
	"github.com/srikanta0427/rest_api_design/model"
	util "github.com/srikanta0427/rest_api_design/util/security"
)

type UserService interface {
	CreateUser(u model.User) (int64, error)
}
type UserServiceImpl struct {
	userRepository db.UserRepository
}

// NewUserService constructor for UserService
func NewUserService(_userRepository db.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: _userRepository,
	}
}

func (u *UserServiceImpl) CreateUser(us model.User) (int64, error) {

	// hashing the password
	password, err2 := util.HashPassWord(us.Password)
	if err2 != nil {
		return 0, err2
	}

	// assign the result to original User
	us.Password = password

	res, err := u.userRepository.Create(&us)
	if err != nil {
		return 0, err
	}
	fmt.Println("Creating user from Service")
	return res, nil
}
