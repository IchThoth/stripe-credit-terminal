package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ichthoth/stripe-credit-terminal/internal/cards"
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
	paymentintent := r.Form.Get("payment_intent")
	paymentcurrency := r.Form.Get("payment_currency")
	paymentamount := r.Form.Get("payment_amount")
	paymentmethod := r.Form.Get("payment_method")

	card := cards.Card{
		Secret: app.config.stripeInfo.secret,
		Key:    app.config.stripeInfo.key,
	}
	pi, err := card.RetrievePaymentIntent(paymentintent)
	if err != nil {
		app.errorLog.Println(err)
	}

	pm, err := card.Getpaymentmethod(paymentmethod)
	if err != nil {
		app.errorLog.Println(err)
	}

	lastFour := pm.Card.Last4
	cardExpiryM := pm.Card.ExpMonth
	cardExpiryY := pm.Card.ExpYear

	data := make(map[string]interface{})
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = paymentintent
	data["pc"] = paymentcurrency
	data["pa"] = paymentamount
	data["pm"] = paymentmethod
	data["last_four"] = lastFour
	data["expiry_month"] = cardExpiryM
	data["expiry_year"] = cardExpiryY
	data["bank_return_code"] = pi.Charges.Data[0].ID

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
