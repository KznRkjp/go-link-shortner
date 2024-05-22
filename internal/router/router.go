package router

import (
	"github.com/KznRkjp/go-link-shortner.git/internal/app"
	"github.com/KznRkjp/go-link-shortner.git/internal/database"
	"github.com/KznRkjp/go-link-shortner.git/internal/middleware/gzipper"
	"github.com/KznRkjp/go-link-shortner.git/internal/middleware/middlelogger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func Main() chi.Router {
	r := chi.NewRouter()
	r.Use(middlelogger.WithLogging)
	r.Use(middleware.Compress(5))
	r.Post("/", gzipper.GzipMiddleware(app.GetURL))
	r.Get("/{id}", gzipper.GzipMiddleware(app.ReturnURL))
	r.Post("/api/shorten", gzipper.GzipMiddleware(app.APIGetURL))
	r.Post("/api/shorten/batch", gzipper.GzipMiddleware(app.APIBatchGetURL))
	r.Get("/ping", gzipper.GzipMiddleware(database.Ping))
	return r
}
