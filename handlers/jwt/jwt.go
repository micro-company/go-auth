package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"io/ioutil"
	"crypto/rsa"
	"encoding/json"
	"github.com/batazor/go-auth/models"
	"time"
)

const (
	PRIVATE_KEY = "cert/private_key.pem"
	PUBLIC_KEY = "cert/public_key.pub"
)

var (
	log       = logrus.New()

	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func initCert() {
	signBytes, err := ioutil.ReadFile(PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
		return
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
		return
	}

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

// Routes creates a REST router
func Routes() chi.Router {
	initCert()

	r := chi.NewRouter()

	r.Get("/debug/:token", Debug)
	r.Post("/", Login)

	return r
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"success\": false}"))
		return
	}

	var user models.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"success\": false}"))
		return
	}

	// Create JWT token
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"success\": false}"))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(tokenString))
	return
}

func Debug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"success\": false}"))
	return
}
