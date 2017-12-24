package middleware

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/micro-company/go-auth/utils"
	"github.com/micro-company/go-auth/utils/recaptcha"
)

func Captcha(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get configuration =======================================================
		ENABLE_CAPTCHA := utils.Getenv("ENABLE_CAPTCHA", "false")
		if ENABLE_CAPTCHA != "true" {
			next.ServeHTTP(w, r)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			utils.Error(w, errors.New(`"`+err.Error()+`"`))
			return
		}

		err = recaptcha.VerifyCaptcha(b)
		if err != nil {
			utils.Error(w, errors.New(`{"captcha":`+err.Error()+`}`))
			return
		}

		// And now set a new body, which will simulate the same data we read:
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

		next.ServeHTTP(w, r)
		return
	})
}
