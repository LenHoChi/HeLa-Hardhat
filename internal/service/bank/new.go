package bank

import "hela-bank-sc/internal/repository/transaction"

type impl struct {
	txRepo transaction.Repository
}

func New(txRepo transaction.Repository) Service {
	return impl{
		txRepo: txRepo,
	}
}
