package payment

import "net/http"

type PaymentHandlers interface {
	HandlePayments(http.ResponseWriter, *http.Request)
	HandleSummary(http.ResponseWriter, *http.Request)
	HandleHealth(http.ResponseWriter, *http.Request)
}
