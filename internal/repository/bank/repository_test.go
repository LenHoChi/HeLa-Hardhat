package bank

import (
	"context"
	"testing"
	"time"

	"hela-bank-sc/internal/domain"

	"github.com/stretchr/testify/require"
)

type historyRow struct {
	address   string
	action    string
	amount    string
	txHash    string
	status    string
	createdAt time.Time
}

func TestRepositoryCreate(t *testing.T) {
	type createInput struct {
		address string
		action  string
		amount  string
		txHash  string
		status  string
	}

	type createExpectedRow struct {
		address string
		action  string
		amount  string
		txHash  string
		status  string
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		input   createInput
		wantErr string
		wantRow *createExpectedRow
	}{
		{
			name: "success",
			input: createInput{
				address: testAddress("create-success"),
				action:  "deposit",
				amount:  "1000000000000000000",
				txHash:  "0xhash-create-success",
				status:  "submitted",
			},
			wantRow: &createExpectedRow{
				address: testAddress("create-success"),
				action:  "deposit",
				amount:  "1000000000000000000",
				txHash:  "0xhash-create-success",
				status:  "submitted",
			},
		},
		{
			name: "invalid amount",
			input: createInput{
				address: testAddress("create-invalid"),
				action:  "deposit",
				amount:  "invalid",
				txHash:  "0xhash-create-invalid",
				status:  "submitted",
			},
			wantErr: "invalid amount: invalid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tx := integrationTx(t)
			repo := New(tx)

			got, err := repo.Create(ctx, tc.input.address, tc.input.action, tc.input.amount, tc.input.txHash, tc.input.status)
			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

			if tc.wantRow == nil {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, tc.wantRow.address, got.Address)
			require.Equal(t, tc.wantRow.action, got.Action)
			require.Equal(t, tc.wantRow.amount, got.Amount)
			require.Equal(t, tc.wantRow.txHash, got.TxHash)
			require.Equal(t, tc.wantRow.status, got.Status)
		})
	}
}

func TestRepositoryListByAddress(t *testing.T) {
	type listByAddressMockConfig struct {
		address     string
		fixtureFile string
	}

	ctx := context.Background()
	tests := []struct {
		name      string
		mock      listByAddressMockConfig
		wantItems []*domain.History
	}{
		{
			name: "success ordered by created_at desc",
			mock: listByAddressMockConfig{
				address:     testAddress("list-success"),
				fixtureFile: "list_by_address_success.sql",
			},
			wantItems: []*domain.History{
				{
					Address:   testAddress("list-success"),
					Action:    "withdraw",
					Amount:    "2000000000000000000",
					TxHash:    "0xnewer",
					Status:    "success",
					CreatedAt: time.Date(2026, 4, 21, 10, 1, 0, 0, time.UTC),
				},
				{
					Address:   testAddress("list-success"),
					Action:    "deposit",
					Amount:    "1000000000000000000",
					TxHash:    "0xolder",
					Status:    "submitted",
					CreatedAt: time.Date(2026, 4, 21, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "empty result",
			mock: listByAddressMockConfig{
				address:     testAddress("list-empty"),
				fixtureFile: "",
			},
			wantItems: []*domain.History{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tx := integrationTx(t)
			repo := New(tx)

			if tc.mock.fixtureFile != "" {
				loadSQLFixture(t, tx, tc.mock.fixtureFile)
			}

			got, err := repo.ListByAddress(ctx, tc.mock.address)
			require.NoError(t, err)
			require.Len(t, got, len(tc.wantItems))

			for i := range tc.wantItems {
				require.Equal(t, tc.wantItems[i].Address, got[i].Address)
				require.Equal(t, tc.wantItems[i].Action, got[i].Action)
				require.Equal(t, tc.wantItems[i].Amount, got[i].Amount)
				require.Equal(t, tc.wantItems[i].TxHash, got[i].TxHash)
				require.Equal(t, tc.wantItems[i].Status, got[i].Status)
				require.True(t, tc.wantItems[i].CreatedAt.Equal(got[i].CreatedAt))
			}
		})
	}
}
