package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router interface {
	Routes() http.Handler
}

type Application struct {
	productRouter Router
}

func (app *Application) Mount() *chi.Mux {
	mux := chi.NewMux()
	mux.Use(EnsureValidToken)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	// mux.Use(middleware.Logger)
	mux.Use(ZapLogger)
	mux.Use(middleware.Recoverer)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Mount("/api/v1/product", app.productRouter.Routes())
	return mux
}

func New(productRouter Router) *Application {
	return &Application{
		productRouter: productRouter,
	}
}
