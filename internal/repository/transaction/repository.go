package transaction

import (
	"context"
	"fmt"
	"time"

	"hela-bank-sc/internal/models"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
)

func (i impl) Create(
	ctx context.Context,
	address string,
	action string,
	amount string,
	txHash string,
	status string,
) error {
	amountDecimal, ok := new(decimal.Big).SetString(amount)
	if !ok {
		return fmt.Errorf("invalid amount: %s", amount)
	}

	record := &models.TransactionHistory{
		Address:   address,
		Action:    action,
		Amount:    types.NewDecimal(amountDecimal),
		TXHash:    txHash,
		Status:    status,
		CreatedAt: time.Now(),
	}

	return record.Insert(ctx, i.db, boil.Infer())
}
