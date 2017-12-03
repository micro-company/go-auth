package session

import (
	"errors"
	"net/http"

	"github.com/micro-company/go-auth/models/session"
	"github.com/micro-company/go-auth/utils"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var TOKEN_ACCESS = r.Header.Get("TOKEN_ACCESS")
		if TOKEN_ACCESS == "" {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New(`"not auth"`))
			return
		}

		token, err := sessionModel.VerifyToken(TOKEN_ACCESS)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New(`"`+err.Error()+`"`))
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			utils.Error(w, errors.New(`"token invalid"`))
		}
		return
	})
}
