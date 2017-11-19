package utils

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	PUBLIC_KEY = "cert/public_key.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	// JWT =====================================================================
	verifyBytes, err := ioutil.ReadFile(PUBLIC_KEY)
	if err != nil {
		log.Fatal(err)
		return
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
		return
	}
}

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
			Error(w, errors.New("not auth"))
			return
		}

		token, err := VerifyToken(TOKEN_ACCESS)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			Error(w, err)
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			Error(w, errors.New("token invalid"))
		}
		return
	})
}
