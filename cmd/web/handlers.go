package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) PosTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplates(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}

	cardHolder := r.Form.Get("cardholder_name")
	email := r.Form.Get("email")
	pi := r.Form.Get("payment_intent")
	pc := r.Form.Get("payment_currency")
	pa := r.Form.Get("payment_amount")
	pm := r.Form.Get("payment_method")

	data := make(map[string]interface{})
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = pi
	data["pc"] = pc
	data["pa"] = pa
	data["pm"] = pm

	if err := app.renderTemplates(w, r, "suceeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	imageId, _ := strconv.Atoi(id)

	image, err := app.DB.GetGopherImages(imageId)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["images"] = image
	if err := app.renderTemplates(w, r, "buy", &templateData{
		Data: data,
	}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}
