package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/srikanta0427/rest_api_design/model"
)

type UserRepository interface {
	Create(u *model.User) (int64, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &UserRepositoryImpl{
		db: db,
	}
}

func (repo *UserRepositoryImpl) Create(u *model.User) (int64, error) {
	fmt.Println("Creating user from Repository")

	// sql query for inserting
	result, err := repo.db.Exec(
		"INSERT INTO user (id, name, email, password, createdAt) VALUES (?, ?, ?, ?,?)",
		u.Id, u.Name, u.Email, u.Password, u.CreatedAt,
	)

	if err != nil {
		return 0, fmt.Errorf("error creating user: %v", err)
	}

	id, err := result.LastInsertId()
	fmt.Println(id)
	return id, err
}
