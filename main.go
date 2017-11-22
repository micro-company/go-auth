package main

import (
	"net/http"
	"time"

	"github.com/batazor/go-auth/db"
	"github.com/batazor/go-auth/handlers/session"
	"github.com/batazor/go-auth/handlers/user"
	"github.com/batazor/go-auth/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var log = logrus.New()

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)

	// Connect to DB
	db.ConnectToMongo()
	db.ConnectToRedis()
}

func main() {

	// Get configuration =======================================================
	PORT := utils.Getenv("PORT", "4070")

	// OpenTracing =============================================================
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  "localhost:5775",
		},
	}
	tracer, closer, _ := cfg.New(
		"go-auth",
		config.Logger(jaeger.StdLogger),
	)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// Routes ==================================================================
	r := chi.NewRouter()

	// CORS ====================================================================
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(cors.Handler)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(utils.NewStructuredLogger(log))
	r.Use(middleware.Recoverer)

	r.Mount("/users", user.Routes())
	r.Mount("/auth", session.Routes())

	// start HTTP-server
	log.Info("Run services on port " + PORT)
	http.ListenAndServe(":"+PORT, r)
}
