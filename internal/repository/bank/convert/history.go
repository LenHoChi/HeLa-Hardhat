package convert

import (
	"hela-bank-sc/internal/domain"
	"hela-bank-sc/internal/models"
)

func ToHistory(record *models.TransactionHistory) *domain.History {
	return &domain.History{
		Address:   record.Address,
		Action:    record.Action,
		Amount:    record.Amount.String(),
		TxHash:    record.TXHash,
		Status:    record.Status,
		CreatedAt: record.CreatedAt,
	}
}
