package router

import (
	"github.com/KznRkjp/go-link-shortner.git/internal/app"
	"github.com/KznRkjp/go-link-shortner.git/internal/middleware/shortlogger"
	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	// "github.com/KznRkjp/go-link-shortner.git/internal/flags"
)

func Main() chi.Router {
	r := chi.NewRouter()
	r.Use(shortlogger.WithLogging)
	r.Post("/", app.GetURL)
	r.Get("/{id}", app.ReturnURL)
	return r
}
