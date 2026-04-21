package bank

import (
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func integrationDB(t *testing.T) *sql.DB {
	t.Helper()

	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.Ping(); err != nil {
		t.Skipf("ping database: %v", err)
	}

	return db
}

func cleanupTransactionHistories(t *testing.T, db *sql.DB, address string) {
	t.Helper()

	_, err := db.Exec(`DELETE FROM transaction_histories WHERE address = $1`, address)
	require.NoError(t, err)
}

func insertHistoryRow(t *testing.T, db *sql.DB, row historyRow) {
	t.Helper()

	_, err := db.Exec(
		`INSERT INTO transaction_histories (address, action, amount, tx_hash, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		row.address,
		row.action,
		row.amount,
		row.txHash,
		row.status,
		row.createdAt,
	)
	require.NoError(t, err)
}

func testAddress(suffix string) string {
	return "itest_" + suffix
}
