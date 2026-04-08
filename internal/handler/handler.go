package handler

import "hela-bank-sc/internal/service/bank"

type BankHandler struct {
	bankSvc bank.Service
}

func NewBankHandler(banksvc bank.Service) BankHandler {
	return BankHandler{bankSvc: banksvc}
}
