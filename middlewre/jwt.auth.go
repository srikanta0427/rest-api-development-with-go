package middlewere

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/srikanta0427/rest_api_design/config"
	"github.com/srikanta0427/rest_api_design/util"
)

func GenerateToken(organizerID int64) (string, error) {
	secret := config.GetString("JWT_SECRET", "")
	if len(secret) < 32 {
		return "", fmt.Errorf("JWT_SECRET must contain at least 32 bytes")
	}

	now := time.Now()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  strconv.FormatInt(organizerID, 10),
			"role": "organizer",
			"iat":  now.Unix(),
			"exp":  now.Add(time.Hour).Unix(),
		},
	)

	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString string) (*jwt.Token, int64, error) {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return nil, 0, fmt.Errorf("JWT_SECRET must contain at least 32 bytes")
	}

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, 0, fmt.Errorf("invalid token")
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return nil, 0, fmt.Errorf("invalid token: %w", err)
	}

	parseInt, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return nil, 0, err
	}

	return token, parseInt, nil
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			util.Error(w, http.StatusBadRequest, "token not found", "")
			return
		}

		verifiedToken, organizerID, err := VerifyToken(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if verifiedToken == nil || !verifiedToken.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "organizerID", organizerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
