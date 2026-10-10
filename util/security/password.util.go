package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassWord(password string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password", err.Error())
	}
	return string(hashPass), nil
}

func ComparePassWord(password, hash string) bool {
	fmt.Println("compare password")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println(err)

		return false
	}
	fmt.Println("hash")
	return true
}
