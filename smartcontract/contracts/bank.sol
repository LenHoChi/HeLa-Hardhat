// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleBank {

    mapping(address => uint) public balances; // address = key, unit = value
    address[] public depositors;

    address public owner;

    error NotOwner();
    error InvalidAmount();
    error InsufficientBalance();
    error TransferFailed();

    event Deposited(address indexed user, uint amount); // logs - see in tab logs in heLa explorer
    event Withdrawn(address indexed user, uint amount);
    event EmergencyWithdrawn(address indexed owner, uint amount);

    constructor() { // owner + contructor => when deploy, owner = sender
        owner = msg.sender;
    }

    modifier onlyOwner() {
        // require(msg.sender==owner, "Not owner");
        if (msg.sender!=owner) revert NotOwner();
        _;
    }

    // deposit contract
    function deposit() public payable {// payable - allow receive money
        if(balances[msg.sender] == 0) {
            depositors.push(msg.sender);  // ← lưu address
        }
        // require(msg.value > 0, "Amount must be greater than 0");
        if (msg.value <= 0) revert InvalidAmount();
        balances[msg.sender] += msg.value; // map [sender] = value / sender = address - calller
    
        emit Deposited(msg.sender, msg.value);
    }

    // withdraw
    function withdraw(uint amount) public {
        // require(amount > 0, "Amount must be greater than 0");
        if(amount <= 0) revert InvalidAmount();
        // require(balances[msg.sender] >= amount, "Not enough balance");
        if(balances[msg.sender] < amount) revert InsufficientBalance();

        balances[msg.sender] -= amount;
        /*
        	•	Trừ số tiền user trong mapping balances
	        •	msg.sender = địa chỉ gọi hàm (user) -> la dia chi ví
        */
        // payable(msg.sender).transfer(amount); // send amount from contract to sender
        (bool success, ) = msg.sender.call{value: amount}("");
        /*
            •	Gửi ETH từ contract → user
        */
        // require(success, "Transfer failed");
        if(!success) revert TransferFailed();

        emit Withdrawn(msg.sender, amount);
    }

    // get all 
    function emergencyWithdraw() public onlyOwner {
        uint balance = address(this).balance;

        // ← reset toàn bộ mapping
        for(uint i = 0; i < depositors.length; i++) {
            balances[depositors[i]] = 0;
        }
        delete depositors;  // ← xóa danh sách

        (bool success, ) = msg.sender.call{value: balance}(""); // send all hlu to meta mask
        // require(success, "Transfer failed");
        if(!success) revert TransferFailed();

        emit EmergencyWithdrawn(msg.sender, balance);
    }

    // Get the balance of the caller stored in contract (internal mapping)
    function getMyBalance() public view returns (uint) {
        return balances[msg.sender];
    }

    // Get the balance of a specific user stored in contract (internal mapping)
    function getBalance(address user) public view returns (uint) {
        return balances[user];
    }

    // Get the actual ETH balance held by this contract (real on-chain balance)
    function getContractBalance() public view returns (uint) {
        return address(this).balance;
    }

    // receive() external payable {
    //     revert("Use deposit()");
    // }

    // Handle direct ETH transfers to contract (without calling deposit())
    // Update internal balance mapping and emit event
    receive() external payable {
        if (msg.value == 0) revert InvalidAmount();

        balances[msg.sender] += msg.value;

        emit Deposited(msg.sender, msg.value);
    }
}