package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/batazor/go-auth/models/user"
	"github.com/batazor/go-auth/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

const (
	PRIVATE_KEY = "cert/private_key.pem"
	PUBLIC_KEY  = "cert/public_key.pub"
)

var (
	log = logrus.New()

	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)

	// JWT =====================================================================
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
		utils.Error(w, err)
		return
	}

	var user userModel.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		utils.Error(w, err)
		return
	}

	var passwordUser = user.Password
	var searchUser = userModel.User{}
	searchUser.Mail = user.Mail
	err, user = userModel.FindOne(searchUser)
	if err != nil {
		utils.Error(w, errors.New("incorrect mail or password"))
		return
	}

	isErr := utils.CheckPasswordHash(passwordUser, user.Password)
	if !isErr {
		utils.Error(w, errors.New("incorrect mail or pass"))
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
		utils.Error(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.Write([]byte("{\"success\": true}"))
	return
}

func Debug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"success\": false}"))
	return
}
