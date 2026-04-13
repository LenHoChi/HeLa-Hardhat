package bank

import (
	"math"

	"github.com/ethereum/go-ethereum/common"
)

func isValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

func isValidAmount(amount float64) bool {
	return amount > 0 && !math.IsNaN(amount) && !math.IsInf(amount, 0)
}
