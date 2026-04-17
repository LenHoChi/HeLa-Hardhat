package bank

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"hela-bank-sc/internal/domain"
	"hela-bank-sc/internal/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestServiceGetBalance(t *testing.T) {
	type getBalanceMockConfig struct {
		balance *big.Int
		err     error
	}

	ctx := context.Background()
	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	balance := big.NewInt(123456789)

	tests := []struct {
		name        string
		mock        getBalanceMockConfig
		wantBalance *big.Int
		wantErr     string
	}{
		{
			name: "success",
			mock: getBalanceMockConfig{
				balance: balance,
			},
			wantBalance: balance,
		},
		{
			name: "chain error",
			mock: getBalanceMockConfig{
				err: errors.New("rpc failed"),
			},
			wantErr: "get balance from blockchain: rpc failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			chain := mocks.NewGateway(t)
			chain.On("GetBalance", ctx, addr).Return(tc.mock.balance, tc.mock.err)

			svc := New(repo, chain)
			got, err := svc.GetBalance(ctx, addr)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wantBalance, got)
		})
	}
}

func TestServiceDeposit(t *testing.T) {
	type depositMockConfig struct {
		txHash      common.Hash
		amountWei   *big.Int
		fromAddress string
		chainErr    error
		repoErr     error
	}

	ctx := context.Background()
	txHash := common.HexToHash("0x123")
	amount := 1.5
	amountWei := big.NewInt(1500000000000000000)
	fromAddress := "0xabc"

	tests := []struct {
		name       string
		mock       depositMockConfig
		wantTxHash common.Hash
		wantErr    string
	}{
		{
			name: "success",
			mock: depositMockConfig{
				txHash:      txHash,
				amountWei:   amountWei,
				fromAddress: fromAddress,
			},
			wantTxHash: txHash,
		},
		{
			name: "chain error",
			mock: depositMockConfig{
				txHash:    txHash,
				amountWei: amountWei,
				chainErr:  errors.New("rpc failed"),
			},
			wantErr: "submit deposit transaction: rpc failed",
		},
		{
			name: "repository error",
			mock: depositMockConfig{
				txHash:      txHash,
				amountWei:   amountWei,
				fromAddress: fromAddress,
				repoErr:     errors.New("db failed"),
			},
			wantErr: "create deposit history: db failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			chain := mocks.NewGateway(t)
			chain.On("Deposit", ctx, amount).Return(tc.mock.txHash, tc.mock.amountWei, tc.mock.chainErr)

			if tc.mock.chainErr == nil {
				chain.On("FromAddress").Return(tc.mock.fromAddress)
				repo.On("Create", ctx, tc.mock.fromAddress, "deposit", tc.mock.amountWei.String(), tc.mock.txHash.Hex(), "submitted").
					Return(tc.mock.repoErr)
			}

			svc := New(repo, chain)
			got, err := svc.Deposit(ctx, amount)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				require.Equal(t, common.Hash{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wantTxHash, got)
		})
	}
}

func TestServiceWithdraw(t *testing.T) {
	type withdrawMockConfig struct {
		txHash      common.Hash
		amountWei   *big.Int
		fromAddress string
		chainErr    error
		repoErr     error
	}

	ctx := context.Background()
	txHash := common.HexToHash("0x456")
	amount := 2.25
	amountWei := big.NewInt(2250000000000000000)
	fromAddress := "0xabc"

	tests := []struct {
		name       string
		mock       withdrawMockConfig
		wantTxHash common.Hash
		wantErr    string
	}{
		{
			name: "success",
			mock: withdrawMockConfig{
				txHash:      txHash,
				amountWei:   amountWei,
				fromAddress: fromAddress,
			},
			wantTxHash: txHash,
		},
		{
			name: "chain error",
			mock: withdrawMockConfig{
				txHash:    txHash,
				amountWei: amountWei,
				chainErr:  errors.New("rpc failed"),
			},
			wantErr: "submit withdraw transaction: rpc failed",
		},
		{
			name: "repository error",
			mock: withdrawMockConfig{
				txHash:      txHash,
				amountWei:   amountWei,
				fromAddress: fromAddress,
				repoErr:     errors.New("db failed"),
			},
			wantErr: "create withdraw history: db failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			chain := mocks.NewGateway(t)
			chain.On("Withdraw", ctx, amount).Return(tc.mock.txHash, tc.mock.amountWei, tc.mock.chainErr)

			if tc.mock.chainErr == nil {
				chain.On("FromAddress").Return(tc.mock.fromAddress)
				repo.On("Create", ctx, tc.mock.fromAddress, "withdraw", tc.mock.amountWei.String(), tc.mock.txHash.Hex(), "submitted").
					Return(tc.mock.repoErr)
			}

			svc := New(repo, chain)
			got, err := svc.Withdraw(ctx, amount)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				require.Equal(t, common.Hash{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wantTxHash, got)
		})
	}
}

func TestServiceEmergencyWithdraw(t *testing.T) {
	type emergencyWithdrawMockConfig struct {
		txHash          common.Hash
		contractBalance *big.Int
		fromAddress     string
		balanceErr      error
		withdrawErr     error
		repoErr         error
	}

	ctx := context.Background()
	txHash := common.HexToHash("0x789")
	contractBalance := big.NewInt(3000000000000000000)
	fromAddress := "0xabc"

	tests := []struct {
		name       string
		mock       emergencyWithdrawMockConfig
		wantTxHash common.Hash
		wantErr    string
	}{
		{
			name: "success",
			mock: emergencyWithdrawMockConfig{
				txHash:          txHash,
				contractBalance: contractBalance,
				fromAddress:     fromAddress,
			},
			wantTxHash: txHash,
		},
		{
			name: "get contract balance error",
			mock: emergencyWithdrawMockConfig{
				txHash:          txHash,
				contractBalance: contractBalance,
				balanceErr:      errors.New("rpc failed"),
			},
			wantErr: "get contract balance before emergency withdraw: rpc failed",
		},
		{
			name: "emergency withdraw error",
			mock: emergencyWithdrawMockConfig{
				txHash:          txHash,
				contractBalance: contractBalance,
				withdrawErr:     errors.New("tx failed"),
			},
			wantErr: "submit emergency withdraw transaction: tx failed",
		},
		{
			name: "repository error",
			mock: emergencyWithdrawMockConfig{
				txHash:          txHash,
				contractBalance: contractBalance,
				fromAddress:     fromAddress,
				repoErr:         errors.New("db failed"),
			},
			wantErr: "create emergency withdraw history: db failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			chain := mocks.NewGateway(t)
			chain.On("GetContractBalance", ctx).Return(tc.mock.contractBalance, tc.mock.balanceErr)

			if tc.mock.balanceErr == nil {
				chain.On("EmergencyWithdraw", ctx).Return(tc.mock.txHash, tc.mock.withdrawErr)
			}

			if tc.mock.balanceErr == nil && tc.mock.withdrawErr == nil {
				chain.On("FromAddress").Return(tc.mock.fromAddress)
				repo.On("Create", ctx, tc.mock.fromAddress, "emergency_withdraw", tc.mock.contractBalance.String(), tc.mock.txHash.Hex(), "submitted").
					Return(tc.mock.repoErr)
			}

			svc := New(repo, chain)
			got, err := svc.EmergencyWithdraw(ctx)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				require.Equal(t, common.Hash{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wantTxHash, got)
		})
	}
}

func TestServiceGetContractBalance(t *testing.T) {
	type getContractBalanceMockConfig struct {
		balance *big.Int
		err     error
	}

	ctx := context.Background()
	balance := big.NewInt(999999)

	tests := []struct {
		name        string
		mock        getContractBalanceMockConfig
		wantBalance *big.Int
		wantErr     string
	}{
		{
			name: "success",
			mock: getContractBalanceMockConfig{
				balance: balance,
			},
			wantBalance: balance,
		},
		{
			name: "chain error",
			mock: getContractBalanceMockConfig{
				err: errors.New("rpc failed"),
			},
			wantErr: "get contract balance from blockchain: rpc failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			chain := mocks.NewGateway(t)
			chain.On("GetContractBalance", ctx).Return(tc.mock.balance, tc.mock.err)

			svc := New(repo, chain)
			got, err := svc.GetContractBalance(ctx)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wantBalance, got)
		})
	}
}

func TestServiceGetHistory(t *testing.T) {
	type getHistoryMockConfig struct {
		histories []*domain.History
		err       error
	}

	ctx := context.Background()
	address := "0x1234567890123456789012345678901234567890"
	createdAt := time.Date(2026, 4, 17, 10, 0, 0, 0, time.UTC)
	histories := []*domain.History{
		{
			Address:   address,
			Action:    "deposit",
			Amount:    "1000000000000000000",
			TxHash:    "0xabc",
			Status:    "submitted",
			CreatedAt: createdAt,
		},
	}

	tests := []struct {
		name          string
		mock          getHistoryMockConfig
		wantHistories []*domain.History
		wantErr       string
	}{
		{
			name: "success",
			mock: getHistoryMockConfig{
				histories: histories,
			},
			wantHistories: histories,
		},
		{
			name: "repository error",
			mock: getHistoryMockConfig{
				err: errors.New("db failed"),
			},
			wantErr: "list history by address: db failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			chain := mocks.NewGateway(t)
			repo.On("ListByAddress", ctx, address).Return(tc.mock.histories, tc.mock.err)

			svc := New(repo, chain)
			got, err := svc.GetHistory(ctx, address)

			if tc.wantErr != "" {
				require.EqualError(t, err, tc.wantErr)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wantHistories, got)
		})
	}
}
