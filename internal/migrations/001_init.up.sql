-- migrations/001_init.up.sql
CREATE TABLE transaction_histories (
    id BIGSERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    action TEXT NOT NULL,
    amount NUMERIC(78, 0) NOT NULL,
    tx_hash TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
