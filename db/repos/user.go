package db

import (
	"database/sql"
	"fmt"
)

type UserRepository interface {
	Create() error
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{

	}
}
type UserRepositoryImpl struct {
	db *sql.DB
}

func (repo *UserRepositoryImpl) Create() error {
	fmt.Println("Creating user from Repository")
	return nil
}
