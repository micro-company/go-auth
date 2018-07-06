package session

import (
	"github.com/micro-company/go-auth/models/session"
	"time"
)

func CreateJWTToken() (string, string, error) {
	// Create JWT token
	timeTTL := time.Minute * 5
	timeDuration := time.Now().Add(timeTTL).Unix()

	// get access token
	tokenString, err := sessionModel.NewAccessToken(timeDuration)
	if err != nil {
		return "", "", err
	}

	// get refresh token
	refreshToken, err := sessionModel.NewRefreshToken(timeTTL)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshToken, nil
}
