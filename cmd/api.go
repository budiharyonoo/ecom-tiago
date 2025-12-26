package main

import (
	"database/sql"
	repositories "github/budiharyonoo/ecom-tiago/internal/adapters/mysql/sqlc"
	"github/budiharyonoo/ecom-tiago/internal/orders"
	"github/budiharyonoo/ecom-tiago/internal/products"
	env "github/budiharyonoo/ecom-tiago/internal/utils"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
)

type app struct {
	config config
	db     *sql.DB
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *app) mount() http.Handler {
	r := chi.NewRouter()

	// logger setup
	appEnv := env.GetString("APP_ENV", "production")
	appName := env.GetString("APP_NAME", "ecom-tiago")
	appVersion := env.GetString("APP_VERSION", "v1.0.0")

	isLocalhost := env.GetString("APP_ENV", "production") == "local"
	logFormat := httplog.SchemaECS.Concise(isLocalhost)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: logFormat.ReplaceAttr,
	})).With(
		slog.String("app", appName),
		slog.String("version", appVersion),
		slog.String("env", appEnv))
	slog.SetDefault(logger)

	// load middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Request logger
	r.Use(httplog.RequestLogger(logger, &httplog.Options{
		// Level defines the verbosity of the request logs:
		// slog.LevelDebug - log all responses (incl. OPTIONS)
		// slog.LevelInfo  - log responses (excl. OPTIONS)
		// slog.LevelWarn  - log 4xx and 5xx responses only (except for 429)
		// slog.LevelError - log 5xx responses only
		Level: slog.LevelDebug,

		// Set log output to Elastic Common Schema (ECS) format.
		Schema: httplog.SchemaECS,

		// RecoverPanics recovers from panics occurring in the underlying HTTP handlers
		// and middlewares. It returns HTTP 500 unless response status was already set.
		//
		// NOTE: Panics are logged as errors automatically, regardless of this setting.
		RecoverPanics: true,

		// Optionally, filter out some request logs.
		Skip: func(req *http.Request, respStatus int) bool {
			return respStatus == 404 || respStatus == 405
		},

		// Optionally, log selected request/response headers explicitly.
		LogRequestHeaders:  []string{"Origin"},
		LogResponseHeaders: []string{},

		// Optionally, enable logging of request/response body based on custom conditions.
		// Useful for debugging payload issues in development.
		LogRequestBody: func(r *http.Request) bool {
			return isLocalhost
		},
		LogResponseBody: func(r *http.Request) bool {
			return isLocalhost
		},
	}))

	// routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	productService := products.NewService(repositories.New(app.db))
	productHandler := products.NewHandler(productService)
	r.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.List)
		r.Get("/{id}", productHandler.GetById)
	})

	orderService := orders.NewService(repositories.New(app.db), app.db)
	orderHandler := orders.NewHandler(orderService)
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.Store)
	})

	return r
}

func (app *app) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 40,
	}

	log.Println("http server started at", app.config.addr)

	return srv.ListenAndServe()
}
