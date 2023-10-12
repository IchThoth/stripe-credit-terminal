package main

import (
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplates(w, r, "terminal", nil); err != nil {
		app.errorLog.Println(err)
	}
}


func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request)  {
	err:= r.ParseForm()
	if err!= nil {
		app.errorLog.Println(err)
	}

	cardHolder := r.Form.Get("cardholder_name")
	email:= r.Form.Get("email")
	pi:= r.Form.Get("payment_intent")
	pc:= r.Form.Get("payment_currency")
	pa:= r.Form.Get("payment_amount")
	pm:=r.Form.Get("payment_method")

	data:= make(map[string]interface{})
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = pi
	data["pc"] = pc
	data["pa"] = pa
	data["pm"] = pm

	if err := app.renderTemplates(w, r, "succeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}
