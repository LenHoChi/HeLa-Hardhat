package router

import (
	"hela-bank-sc/internal/handler"
	"hela-bank-sc/internal/service/bank"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	bankHandler *handler.BankHandler
}

// func NewRouter() chi.Router {
// 	r := chi.NewRouter()
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)

// 	r.Get("/balance/{address}", GetBalance)
// 	// deposit(1) đang được hiểu là 1 ETH theo đơn vị wei = 10^18
// 	//   =      1,000,000,000,000,000,000
// 	// 1 coin = 1 000 000 000 000 000 000 wei
// 	r.Post("/deposit", Deposit)
// 	r.Post("/withdraw", Withdraw)

// 	return r
// }

func NewRouter(bankSvc bank.Service) chi.Router {
	router := &Router{
		bankHandler: handler.NewBankHandler(bankSvc),
	}

	// other way:
	// handler := 	NewBankHandler(bankSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/balance/{address}", router.bankHandler.GetBalance())
	// r.Get("/balance/{address}", handler.GetBalance())

	// deposit(1) đang được hiểu là 1 ETH theo đơn vị wei = 10^18
	//   =      1,000,000,000,000,000,000
	// 1 coin = 1 000 000 000 000 000 000 wei
	r.Post("/deposit", router.bankHandler.Deposit())
	r.Post("/withdraw", router.bankHandler.Withdraw())

	return r
}
