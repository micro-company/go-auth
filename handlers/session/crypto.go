package session

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/batazor/go-auth/utils"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, err
}

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var TOKEN_ACCESS = r.Header.Get("TOKEN_ACCESS")
		if TOKEN_ACCESS == "" {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New("not auth"))
			return
		}

		token, err := VerifyToken(TOKEN_ACCESS)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, err)
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New("token invalid"))
		}
		return
	})
}
