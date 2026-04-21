package bank

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "valid hex address",
			address: "0x1234567890123456789012345678901234567890",
			want:    true,
		},
		{
			name:    "empty address",
			address: "",
			want:    false,
		},
		{
			name:    "missing 0x prefix",
			address: "1234567890123456789012345678901234567890",
			want:    true,
		},
		{
			name:    "too short",
			address: "0x1234",
			want:    false,
		},
		{
			name:    "non hex characters",
			address: "0xZZZZ567890123456789012345678901234567890",
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, isValidAddress(tc.address))
		})
	}
}

func TestIsValidAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		want   bool
	}{
		{
			name:   "positive amount",
			amount: 1.5,
			want:   true,
		},
		{
			name:   "zero amount",
			amount: 0,
			want:   false,
		},
		{
			name:   "negative amount",
			amount: -1,
			want:   false,
		},
		{
			name:   "nan amount",
			amount: math.NaN(),
			want:   false,
		},
		{
			name:   "positive infinity",
			amount: math.Inf(1),
			want:   false,
		},
		{
			name:   "negative infinity",
			amount: math.Inf(-1),
			want:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, isValidAmount(tc.amount))
		})
	}
}
