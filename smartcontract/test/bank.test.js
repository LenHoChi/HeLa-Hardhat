const { expect } = require("chai");
const hre = require("hardhat");

describe("SimpleBank", function () {
  async function deploy() {
    const [owner, user1] = await hre.ethers.getSigners(); // deploy contract before testing
    const bank = await hre.ethers.deployContract("SimpleBank");
    return { bank, owner, user1 };  
    // bank smart contract after deploy
    // owner - msg.sender - wallet deploy contract
    // user 1 - another wallet, 
  }

  it("should deposit", async function () {
    const { bank, user1 } = await deploy();

    await bank.connect(user1).deposit({ value: hre.ethers.parseEther("1") }); // user1 call deposit with 1 eth

    expect(await bank.getBalance(user1.address)).to.equal( // get balance of user1
      hre.ethers.parseEther("1")
    );
  });

  it("should withdraw", async function () {
    const { bank, user1 } = await deploy();

    await bank.connect(user1).deposit({ value: hre.ethers.parseEther("1") });
    await bank.connect(user1).withdraw(hre.ethers.parseEther("0.5"));

    expect(await bank.getBalance(user1.address)).to.equal(
      hre.ethers.parseEther("0.5")
    );
  });

  it("should emergency withdraw", async function () {
    const { bank, owner, user1 } = await deploy();

    await bank.connect(user1).deposit({ value: hre.ethers.parseEther("1") });
    await bank.connect(owner).emergencyWithdraw();

    expect(await hre.ethers.provider.getBalance(bank.target)).to.equal(0);
  });
});

// muc dich contract 
//User gửi tiền vào contract → contract giữ tiền → user có thể rút lại sau
// Hiểu như ngân hàng mini
// Contract = ngân hàng
//Mapping = sổ tài khoản
/*
🔹 Tại sao không dùng ví luôn?
Ví chỉ là:
	•	giữ tiền cá nhân không có logic
Smart contract có thể:
	•	kiểm soát tiền theo rule | khóa tiền | chia tiền | staking | lending | DAO | DeFi
🔹 Ví dụ thực tế

1. Escrow (trung gian)
	•	A gửi tiền vào contract
	•	chỉ khi B hoàn thành việc → contract mới trả tiền
2. Staking
	•	bạn gửi ETH vào contract
	•	contract giữ → trả reward
3. Game / NFT
	•	tiền vào contract
	•	contract quản lý logic
*/