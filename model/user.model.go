package model

import "time"

type User struct {
	Id        int       `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) UserCons(id int, name, email, password string, createdAt time.Time) User {
	return User{
		Id:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
	}
}
