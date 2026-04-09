package bank

import (
	"hela-bank-sc/internal/blockchain"
	"hela-bank-sc/internal/repository/transaction"
)

type impl struct {
	txRepo transaction.Repository
	chain  blockchain.Gateway
}

func New(txRepo transaction.Repository, chain blockchain.Gateway) Service {
	return impl{
		txRepo: txRepo,
		chain:  chain,
	}
}
