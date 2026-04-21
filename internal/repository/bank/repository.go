package bank

import (
	"context"
	"fmt"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"

	"hela-bank-sc/internal/domain"
	"hela-bank-sc/internal/models"
	"hela-bank-sc/internal/repository/bank/convert"
)

func (i impl) Create(
	ctx context.Context,
	address string,
	action string,
	amount string,
	txHash string,
	status string,
) (*domain.History, error) {
	amountDecimal, ok := new(decimal.Big).SetString(amount)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", amount)
	}

	record := &models.TransactionHistory{
		Address:   address,
		Action:    action,
		Amount:    types.NewDecimal(amountDecimal),
		TXHash:    txHash,
		Status:    status,
		CreatedAt: time.Now(),
	}

	if err := record.Insert(ctx, i.db, boil.Infer()); err != nil {
		return nil, err
	}

	return convert.ToHistory(record), nil
}

func (i impl) ListByAddress(ctx context.Context, address string) ([]*domain.History, error) {
	records, err := models.TransactionHistories(
		models.TransactionHistoryWhere.Address.EQ(address),
		qm.OrderBy("created_at desc"),
	).All(ctx, i.db)
	if err != nil {
		return nil, err
	}

	items := make([]*domain.History, 0, len(records))
	for _, record := range records {
		items = append(items, convert.ToHistory(record))
	}

	return items, nil
}
