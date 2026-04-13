package domain

import "time"

type History struct {
	Address   string
	Action    string
	Amount    string
	TxHash    string
	Status    string
	CreatedAt time.Time
}
