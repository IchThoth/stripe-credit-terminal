package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/home", app.Home)
	mux.Post("/payment-succeded",app.PaymentSucceeded)

	return mux
}
