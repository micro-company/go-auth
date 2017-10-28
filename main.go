package main

import (
	"github.com/batazor/go-auth/db"
	"github.com/batazor/go-auth/handlers/jwt"
	"github.com/batazor/go-auth/handlers/user"
	"github.com/batazor/go-auth/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"net/http"
	"time"
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
