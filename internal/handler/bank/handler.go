package bank

import "hela-bank-sc/internal/service/bank"

type Handler struct {
	bankSvc bank.Service
}

func New(banksvc bank.Service) Handler {
	return Handler{bankSvc: banksvc}
}
