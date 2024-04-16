package router

import (
	"github.com/KznRkjp/go-link-shortner.git/internal/app"
	"github.com/go-chi/chi/v5"
	// "github.com/KznRkjp/go-link-shortner.git/internal/flags"
)

func Main() chi.Router {
	r := chi.NewRouter()
	r.Post("/", app.GetURL)
	r.Get("/{id}", app.ReturnURL)
	return r
}
