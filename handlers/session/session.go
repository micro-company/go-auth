package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/micro-company/go-auth/utils/crypto"

	"github.com/micro-company/go-auth/utils/recaptcha"

	"github.com/go-chi/chi"
	"github.com/micro-company/go-auth/handlers/user"
	"github.com/micro-company/go-auth/models/session"
	"github.com/micro-company/go-auth/models/user"
	"github.com/micro-company/go-auth/utils"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

// Routes creates a REST router
func Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/debug/:token", Debug)
	r.Post("/", Login)
	r.Post("/new", Registration)
	r.Post("/recovery", Recovery)
	r.Delete("/", Logout)

	return r
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	var user userModel.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	// Check recaptcha
	err = recaptcha.VerifyCaptcha(b)
	if err != nil {
		utils.Error(w, errors.New(`{"captcha":`+err.Error()+`}`))
		return
	}

	var passwordUser = user.Password
	var searchUser = userModel.User{}
	searchUser.Mail = user.Mail
	err, user = userModel.FindOne(searchUser)
	if err != nil {
		utils.Error(w, errors.New(`{"mail":"incorrect mail or password"}`))
		return
	}

	isErr := crypto.CheckPasswordHash(passwordUser, user.Password)
	if !isErr {
		utils.Error(w, errors.New(`{"mail":"incorrect mail or password"}`))
		return
	}

	// Create JWT token
	timeDuration := time.Now().Add(time.Minute * 5).Unix()

	// get access token
	tokenString, err := sessionModel.NewAccessToken(timeDuration)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	// get refresh token
	refreshToken, err := sessionModel.NewRefreshToken(timeDuration)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{
		"tokens": {
			"access": "` + tokenString + `",
			"refresh": "` + refreshToken + `"
		}
	}`))
}

func Registration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	// Check recaptcha
	err = recaptcha.VerifyCaptcha(b)
	if err != nil {
		utils.Error(w, errors.New(`{"captcha":`+err.Error()+`}`))
		return
	}

	// And now set a new body, which will simulate the same data we read:
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	user.Create(w, r)
}

func Recovery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	// Check recaptcha
	err = recaptcha.VerifyCaptcha(b)
	if err != nil {
		utils.Error(w, errors.New(`{"captcha":`+err.Error()+`}`))
		return
	}

	var user userModel.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	// TODO: Create template mail.tpl
	// TODO: generate link for recovery pass (save URL and TTL to redis)
	// TODO: Load template (use mail.tpl) and render valid html
	// TODO: send mail

	// TODO: FRONT-END
	// TODO: router `/recovery/:id`
	// TODO: write new password
	// TODO: request to BACK-END

	// TODO: BACK-END
	// TODO: search user
	// TODO: id false -> return error
	// TODO:    true -> save new password(hash)
	// TODO:         -> return success
	// TODO: Send mail (theme: New password)
	// TODO: Delete `/recovery/:id` from DB
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var TOKEN_REFRESH = r.Header.Get("TOKEN_REFRESH")
	if TOKEN_REFRESH == "" {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"not auth"`))
		return
	}

	// Chech REFRESH TOKEN
	status, err := sessionModel.CheckRefreshToken(TOKEN_REFRESH)
	if err != nil || status != true {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	// Create JWT token
	timeDuration := time.Now().Add(time.Minute * 5).Unix()

	// get access token
	tokenString, err := sessionModel.NewAccessToken(timeDuration)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	// get refresh token
	refreshToken, err := sessionModel.NewRefreshToken(timeDuration)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{
		"tokens": {
			"access": "` + tokenString + `",
			"refresh": "` + refreshToken + `"
		}
	}`))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Authorization = r.Header.Get("Authorization")
	if Authorization == "" {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"not auth"`))
		return
	}

	token, err := sessionModel.VerifyToken(Authorization)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"token invalid"`))
		return
	}

	err = sessionModel.Delete(Authorization)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func Debug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{}`))
}
