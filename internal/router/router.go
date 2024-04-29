package router

import (
	"github.com/KznRkjp/go-link-shortner.git/internal/app"
	"github.com/KznRkjp/go-link-shortner.git/internal/middleware/middlelogger"
	"github.com/go-chi/chi/v5"
)

func Main() chi.Router {
	r := chi.NewRouter()
	r.Use(middlelogger.WithLogging)
	r.Post("/", app.GetURL)
	r.Get("/{id}", app.ReturnURL)
	r.Post("/api/shorten", app.APIGetURL)
	return r
}
