package payment

import "net/http"

type PaymentHandlers interface {
	CreatePayment(http.ResponseWriter, *http.Request)
}
