package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/pos-terminal", app.PosTerminal)
	mux.Post("/payment-succeded", app.PaymentSucceeded)
	mux.Get("/charge-once", app.ChargeOnce)

	fileserver := http.FileServer(http.Dir("./static"))
	mux.Handle("./static/*", http.StripPrefix("/static", fileserver))
	return mux
}
