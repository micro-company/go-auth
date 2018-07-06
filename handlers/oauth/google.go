package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/micro-company/go-auth/handlers/session"
	"github.com/micro-company/go-auth/models/user"
	"github.com/micro-company/go-auth/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

var (
	// Get configuration
	ClientID     = utils.Getenv("OAUTH_GOOGLE_CLIENT_ID", "YOUR_CLIENT_ID")
	ClientSecret = utils.Getenv("OAUTH_GOOGLE_CLIENT_SECRET", "YOUR_CLIENT_SECRET")
	RedirectURL  = utils.Getenv("OAUTH_REDIRECT_URL", "http://localhost:3000/auth/callback/:type")

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  RedirectURL,
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
)

func googleOAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"url": "` + url + `"}`))
}

func googleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	var oauthCallback Callback
	err = json.Unmarshal(b, &oauthCallback)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), oauthCallback.Code)
	if err != nil {
		utils.Error(w, errors.New("\"cannot fetch token\""))
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	var userGoogle UserGoogle
	err = json.Unmarshal(contents, &userGoogle)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return
	}

	var searchUser = userModel.User{Email: &userGoogle.Email}
	_, err = userModel.FindOne(searchUser)
	if err != nil {
		searchUser.Gender = userGoogle.Gender
		userModel.Add(searchUser)
	}

	// Create JWT token
	tokenString, refreshToken, err := session.CreateJWTToken()
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
