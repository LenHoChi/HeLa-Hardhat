# Hela Bank SC Go

Go backend service for interacting with a Hela testnet smart contract.  
This project exposes REST APIs for deposit, withdraw, emergency withdraw, user balance lookup, contract balance lookup, and transaction history.

## Features

- Connect to Hela testnet RPC
- Interact with a bank smart contract through ABI calls
- Submit blockchain transactions from the backend
- Store transaction history in PostgreSQL
- Listen for contract activity through explorer polling
- Unit tests for handler and service layers
- Integration tests for repository layer

## Project Structure

```text
.
├── abi/                     # Contract ABI
├── internal/
│   ├── blockchain/          # RPC, contract, tx, event logic
│   ├── config/              # Env-based config loading
│   ├── database/            # PostgreSQL connection
│   ├── handler/bank/        # HTTP handlers
│   ├── migrations/          # SQL migrations
│   ├── mocks/               # mockery-generated mocks
│   ├── repository/bank/     # DB access
│   ├── router/              # Route registration
│   └── service/bank/        # Business logic
├── smartcontract/           # Smart contract workspace
├── main.go
└── Makefile
```

## Smart Contract Deployment

### Working Directory
- `./smartcontract`

### Steps

Step 1
- Go to `./smartcontract/contracts`
- Add your Solidity contract file (`.sol`) to this folder

Step 2
- Go to `./smartcontract/scripts`
- Update `deploy.js` to deploy your contract

Step 3
- Run:
  - `npx hardhat compile`
  - `npx hardhat run scripts/deploy.js --network hela_testnet`

Step 4
- Copy the deployed contract address into `CONTRACT_ADDRESS` in `.env`

Step 5
- Copy the latest ABI from `smartcontract/artifacts/contracts/bank.sol/SimpleBank.json`
- Update `abi/bank.json` in the root project

## Prerequisites for Running the Go Backend

- Go `1.24.5`
- Docker / Docker Compose or PostgreSQL
- `golang-migrate` CLI
- `mockery` CLI

Optional but useful:

- VS Code with Go extension

## Environment Variables

Create a `.env` file in the project root.

Example:

```env
HELA_TESTNET_RPC=https://666888.rpc.thirdweb.com
PRIVATE_KEY=your_private_key
DATABASE_URL=postgres://hela-api:hela-api@localhost:5450/hela-bank-sc?sslmode=disable

APP_PORT=8080

WSS_URL=wss://testnet-rpc.helachain.com
CONTRACT_ADDRESS=0x85933342B34ceB2ef5ECc63FEC7659c4a3495d6F

HTTP_READ_HEADER_TIMEOUT_SECONDS=5
HTTP_READ_TIMEOUT_SECONDS=10
HTTP_WRITE_TIMEOUT_SECONDS=10
HTTP_IDLE_TIMEOUT_SECONDS=60

POSTGRES_USER=hela-api
POSTGRES_PASSWORD=hela-api
POSTGRES_DB=hela-bank-sc
```

Required:

- `DATABASE_URL`
- `HELA_TESTNET_RPC`
- `WSS_URL`
- `CONTRACT_ADDRESS`
- `PRIVATE_KEY`

## Install Dependencies

```bash
go mod download
```

## Running the Backend

You can run this backend in two ways.

### Option 1: Manual Setup

Start only PostgreSQL and run migration:

```bash
make setup
```

Rollback one migration:

```bash
make migrate-down
```

If your migration state becomes dirty, fix it with `migrate force` before running again.

Run the application:

```bash
go run main.go
```

The server starts on:

```text
http://localhost:8080
```

or the port configured by `APP_PORT`.

### Option 2: Run with Docker Compose

Build and start PostgreSQL, migration, and the Go app:

```bash
make docker-up
```

Run in detached mode:

```bash
make docker-up-detached
```

Stop all services:

```bash
make docker-down
```

The `docker-compose.yml` file uses these variables from `.env`:

- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`
- `APP_PORT`

The app service also overrides `DATABASE_URL` inside the container so it connects to the `postgres` service instead of `localhost`.

Useful database commands:

```bash
make docker-db-up -> start only PostgreSQL
make setup -> start DB and run migration
make db-clear -> remove DB data
make reset -> clear DB and re-run setup
```

## API Endpoints

### Health

- `GET /check-health`

### Balance

- `GET /balance/{address}`
- `GET /contract-balance`

### Transactions

- `POST /deposit`
- `POST /withdraw`
- `POST /emergency-withdraw`

### History

- `GET /history/{address}`

## Example Requests

### Deposit

```bash
curl --location 'http://localhost:8080/deposit' \
  --header 'Content-Type: application/json' \
  --data '{
    "amount": 1.5
  }'
```

### Withdraw

```bash
curl --location 'http://localhost:8080/withdraw' \
  --header 'Content-Type: application/json' \
  --data '{
    "amount": 0.5
  }'
```

### Get Balance

```bash
curl --location 'http://localhost:8080/balance/0x1234567890123456789012345678901234567890'
```

### Get History

```bash
curl --location 'http://localhost:8080/history/0x1234567890123456789012345678901234567890'
```

## Tests

Run handler tests:

```bash
make test-handler
```

Run service tests:

```bash
make test-service
```

Run repository tests:

```bash
make test-repo
```

Run router tests:

```bash
make test-router
```

Run all main bank-related tests:

```bash
make test-bank
```

Run all tests:

```bash
make test-all
```

Note:

- `test-repo` is an integration test and requires `DATABASE_URL` to be available in the environment.

## Mock Generation

Generate mocks:

```bash
make mock-service
make mock-repo
make mock-gateway
```

Generate all mocks:

```bash
make mock-all
```

Clean and regenerate all mocks:

```bash
make mock-refresh
```

## Debugging

The repository includes VS Code launch configurations for:

- `Run Go App`
- `Debug Go App`
- `Debug Handler Tests`
- `Debug Service Tests`
- `Debug Repository Tests`

If you debug repository tests, make sure `DATABASE_URL` is available through `.env` or your environment.

## Notes

- Repository tests use real PostgreSQL integration tests with transaction-per-test and rollback.
- Handler and service tests use `mockery`-generated mocks.
- This repository is currently the backend part of the project. A frontend/UI can be added later on top of the REST API.
