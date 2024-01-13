package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ichthoth/stripe-credit-terminal/internal/cards"
	"github.com/ichthoth/stripe-credit-terminal/internal/models"
)

func (app *application) PosTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplates(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplates(w, r, "terminal", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}
func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
	}

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	paymentintent := r.Form.Get("payment_intent")
	paymentcurrency := r.Form.Get("payment_currency")
	paymentamount := r.Form.Get("payment_amount")
	paymentmethod := r.Form.Get("payment_method")
	gophimg, _ := strconv.Atoi(r.Form.Get("product_id"))

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

	// creates the customer
	customerID, err := app.SaveCustomer(firstName, lastName, email)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Println(customerID)

	//creates a new transaction
	amount, err := strconv.Atoi(paymentamount)
	if err != nil {
		app.errorLog.Println(err)
	}

	txn := models.Transaction{
		Amount:              amount,
		Currency:            paymentcurrency,
		LastFour:            lastFour,
		BankReturnCode:      pi.Charges.Data[0].ID,
		ExpiryMonth:         int(cardExpiryM),
		ExpiryYear:          int(cardExpiryY),
		TransactionStatusID: 2,
	}

	txnID, err := app.SaveTransaction(txn)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	//creating a new order
	ord := models.Order{
		ImageID:       gophimg,
		TransactionID: txnID,
		CustomerID:    customerID,
		StatusID:      1,
		Quantity:      1,
		Amount:        amount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = app.SaveOrder(ord)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(map[string]interface{})
	data["email"] = email
	data["fn"] = firstName
	data["ln"] = lastName
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

func (app *application) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}

	return id, nil

}

func (app *application) SaveTransaction(txn models.Transaction) (int, error) {
	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (app *application) SaveOrder(ord models.Order) (int, error) {
	id, err := app.DB.InsertOrder(ord)
	if err != nil {
		return 0, err
	}

	return id, nil
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
