package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)

	mux.Post("/authenticate", app.authenticate)
	mux.Post("/register", app.register)
	mux.Get("/refresh", app.refreshToken)

	mux.Get("/products", app.AllProducts)
	mux.Get("/products/{id}", app.GetProductByID)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authRequired)
		mux.Get("/products/{id}/coupon", app.GenerateCoupon)

	})
	mux.Post("/purchase", app.Purchase)

	mux.Get("/history", app.HistoriesByUser)
	//mux.

	return mux
}
