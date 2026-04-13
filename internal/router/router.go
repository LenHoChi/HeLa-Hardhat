package router

import (
	"context"
	handlerbank "hela-bank-sc/internal/handler/bank"
	"hela-bank-sc/internal/service/bank"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	Ctx     context.Context
	BankSvc bank.Service
}

func (rtr Router) Routes() chi.Router {
	// at router: create handler to call svc func (that's why need create svc from main and inject them in struct - will use these svc here)
	// flow: main > router > handler > svc
	handler := handlerbank.New(rtr.BankSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/balance/{address}", handler.GetBalance())

	// 1 coin = 1 000 000 000 000 000 000 wei (10^18)
	r.Post("/deposit", handler.Deposit())

	r.Post("/withdraw", handler.Withdraw())

	r.Post("/emergency-withdraw", handler.EmergencyWithdraw())

	r.Get("/contract-balance", handler.GetContractBalance())

	r.Get("/check-health", handler.CheckHealth())

	r.Get("/history/{address}", handler.GetHistory())

	return r
}
