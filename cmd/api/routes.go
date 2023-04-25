package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", app.Home)

	mux.Get("/authenticate", app.authenticate)
	mux.Get("/products", app.AllProducts)

	return mux
}
