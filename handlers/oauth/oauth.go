package oauth

import (
	"github.com/sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/micro-company/go-auth/middleware"
	"net/http"
	"google.golang.org/appengine"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine/urlfetch"
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

	r.Get("/google", Google)

	return r
}

func Google(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := appengine.NewContext(r)

	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(ctx, "https://www.googleapis.com/auth/bigquery"),
			Base: &urlfetch.Transport{
				Context: ctx,
			},
		},
	}

	client.Get("...")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}