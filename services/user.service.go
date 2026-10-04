package services

import (
	"fmt"

	db "github.com/srikanta0427/rest_api_design/db/repos"
)

type UserService interface {
	CreateUser() error
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

func (u *UserServiceImpl) CreateUser() error {
	_ = u.userRepository.Create()
	fmt.Println("Creating user from Service")
	return nil

}
