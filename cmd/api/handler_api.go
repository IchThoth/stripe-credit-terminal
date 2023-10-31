package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ichthoth/stripe-credit-terminal/internal/cards"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type response struct {
	ID      int    `json:"id,omitempty"`
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		app.errorLog.Println(err)
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
	}

	card := cards.Card{
		Currency: payload.Currency,
		Secret:   app.config.stripeInfo.secret,
		Key:      app.config.stripeInfo.key,
	}

	okay := true

	paymentIntent, msg, err := card.ChargeCard(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		out, err := json.MarshalIndent(paymentIntent, " ", "")
		if err != nil {
			app.errorLog.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		resp := response{
			OK:      false,
			Message: msg,
			Content: "",
		}

		jsonOut, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonOut)
	}

}

func (app *application) GetGopherImagesById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetId, _ := strconv.Atoi(id)

	widget, err := app.DB.GetGopherImages(widgetId)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	jsonOut, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonOut)

}
