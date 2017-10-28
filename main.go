package main

import (
	"github.com/batazor/go-auth/db"
	"github.com/batazor/go-auth/handlers/jwt"
	"github.com/batazor/go-auth/handlers/user"
	"github.com/batazor/go-auth/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.New()

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)

	// Connect to MongoDB
	db.Connect()
}

func main() {

	// Get configuration ======================================================
	PORT := utils.Getenv("PORT", "4070")

	// Routes ==================================================================
	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(utils.NewStructuredLogger(log))
	r.Use(middleware.Recoverer)

	r.Mount("/users", user.Routes())
	r.Mount("/jwt", jwt.Routes())

	// start HTTP-server
	log.Info("Run services on port " + PORT)
	http.ListenAndServe(":"+PORT, r)
}
