package oauth

import (
	"github.com/go-chi/chi"
	"github.com/micro-company/go-auth/middleware"
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
	r.Use(middleware.Captcha)

	r.Get("/google", googleOAuth)
	r.Post("/callback/google", googleCallback)

	return r
}
