package transaction

import "context"

type Repository interface {
	Create(
		ctx context.Context,
		address string,
		action string,
		amount string,
		txHash string,
		status string,
	) error
}
