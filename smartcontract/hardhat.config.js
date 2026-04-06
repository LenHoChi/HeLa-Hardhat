require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config();

module.exports = {
  solidity: "0.8.28",
  networks: {
    hela_testnet: {
      url: process.env.HELA_TESTNET_RPC || "",
      chainId: 666888,
      accounts: [process.env.PRIVATE_KEY || ""],
    },
  },
};