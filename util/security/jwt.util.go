package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/srikanta0427/rest_api_design/config"
)

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.GetString("SECRET_KEY", "my-secret")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
