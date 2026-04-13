package bank

import (
	"hela-bank-sc/internal/blockchain"
	repositorybank "hela-bank-sc/internal/repository/bank"
)

type impl struct {
	txRepo repositorybank.Repository
	chain  blockchain.Gateway
}

func New(txRepo repositorybank.Repository, chain blockchain.Gateway) Service {
	return impl{
		txRepo: txRepo,
		chain:  chain,
	}
}
