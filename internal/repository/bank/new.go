package bank

import (
	"github.com/aarondl/sqlboiler/v4/boil"
)

type impl struct {
	db boil.ContextExecutor
}

func New(db boil.ContextExecutor) Repository {
	return impl{db: db}
}
