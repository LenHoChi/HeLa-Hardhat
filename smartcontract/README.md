# Sample Hardhat Project - Smart Contract Workspace

This folder contains the Hardhat workspace for the bank smart contract used by the Go backend.

## Stack

- Hardhat
- Solidity `0.8.28`
- Hela testnet (`chainId: 666888`)

## Project Structure

```text
smartcontract/
├── contracts/
│   ├── bank.sol
│   └── Lock.sol
├── scripts/
│   └── deploy.js
├── test/
│   ├── bank.test.js
│   ├── Lock.js
│   └── test-demo.js
├── hardhat.config.js
├── package.json
└── .env
```

## Environment Variables

Create a `.env` file inside `smartcontract/`.

Example:

```env
HELA_TESTNET_RPC=https://666888.rpc.thirdweb.com
PRIVATE_KEY=your_private_key
```

Required:

- `HELA_TESTNET_RPC`
- `PRIVATE_KEY`

## Install Dependencies

From the `smartcontract/` directory:

```bash
npm install
```

## Compile Contracts

```bash
npx hardhat compile
```

## Run Tests

Run all tests:

```bash
npx hardhat test
```

Run only the bank contract tests:

```bash
npx hardhat test test/bank.test.js
```

## Deploy to Hela Testnet

The current deploy script deploys the `SimpleBank` contract.

Run:

```bash
npx hardhat run scripts/deploy.js --network hela_testnet
```

After deployment:

1. Copy the deployed contract address from the terminal output.
2. Update `CONTRACT_ADDRESS` in the root project `.env`.
3. Make sure the backend uses the same RPC endpoint and chain configuration.

## Notes

- `contracts/Lock.sol` and some sample test files are still present from the default Hardhat scaffold.
- The active deployment flow currently uses `contracts/bank.sol` through `scripts/deploy.js`.
