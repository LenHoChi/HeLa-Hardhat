package bank

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func integrationDB(t *testing.T) *sql.DB {
	t.Helper()

	// err := godotenv.Load()

	// if err := godotenv.Load("../../../.env", ".env"); err != nil {
	// 	t.Logf("load .env: %v", err)
	// }

	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Logf("load .env: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Logf("load db: %v", err)
		require.NoError(t, err)
		// t.Skip("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Logf("load pq: %v", err)
		require.NoError(t, err)
		// t.Skipf("open database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.Ping(); err != nil {
		t.Skipf("ping database: %v", err)
	}

	return db
}

func integrationTx(t *testing.T) *sql.Tx {
	t.Helper()

	db := integrationDB(t)

	tx, err := db.Begin()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	return tx
}

func loadSQLFixture(t *testing.T, exec boil.ContextExecutor, name string) {
	t.Helper()

	path := filepath.Join("testdata", name)

	sqlBytes, err := os.ReadFile(path)
	require.NoError(t, err)

	_, err = exec.Exec(string(sqlBytes))
	require.NoError(t, err)
}

func testAddress(suffix string) string {
	return "itest_" + suffix
}

// func insertHistoryRow(t *testing.T, exec boil.ContextExecutor, row historyRow) {
// 	t.Helper()

// 	_, err := exec.Exec(
// 		`INSERT INTO transaction_histories (address, action, amount, tx_hash, status, created_at)
// 		 VALUES ($1, $2, $3, $4, $5, $6)`,
// 		row.address,
// 		row.action,
// 		row.amount,
// 		row.txHash,
// 		row.status,
// 		row.createdAt,
// 	)
// 	require.NoError(t, err)
// }

// func cleanupTransactionHistories(t *testing.T, db *sql.DB, address string) {
// 	t.Helper()

// 	_, err := db.Exec(`DELETE FROM transaction_histories WHERE address = $1`, address)
// 	require.NoError(t, err)
// }
