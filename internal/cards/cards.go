package cards

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	//charge card intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMesages(stripeErr.Code)
		}
		return nil, msg, err
	}
	return pi, "", nil
}

func (c *Card) ChargeCard(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

func cardErrorMesages(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card declined please"
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card has expired"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Incorrect CVC"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Incorrect postal code"
	case stripe.ErrorCodeIncorrectZip:
		msg = "Incorrect zip code"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "This amount is too large"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "This amount is too small"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "Insufficient balance"
	default:
		msg = "Your card has declined"
	}

	return msg
}
