package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	log       = logrus.New()
	SecretKey = "WOW,MuchShibe,ToDogge"
)

// Routes creates a REST router
func Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/debug/:token", Debug)
	r.Post("/", Auth)

	return r
}

func Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: check user
	// Create JWT token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	//token.Claims["userid"] = "userId"

	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tokenString))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"success\": false}"))
	return
}

func Debug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"success\": false}"))
	return
}
