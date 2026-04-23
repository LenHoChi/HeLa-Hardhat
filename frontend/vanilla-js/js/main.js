import {
  API_BASE_URL,
  deposit,
  getBalance,
  getHistory,
  withdraw,
} from "./api.js";

const addressInput = document.querySelector("#address-input");
const loadBalanceButton = document.querySelector("#load-balance-btn");
const loadHistoryButton = document.querySelector("#load-history-btn");
const submitDepositButton = document.querySelector("#submit-deposit-btn");
const submitWithdrawButton = document.querySelector("#submit-withdraw-btn");
const depositAmount = document.querySelector("#deposit-amount");
const withdrawAmount = document.querySelector("#withdraw-amount");
const resultAddress = document.querySelector("#result-address");
const resultBalance = document.querySelector("#result-balance");
const historyResult = document.querySelector("#history-result");
const apiBaseUrl = document.querySelector("#api-base-url");

apiBaseUrl.textContent = API_BASE_URL;

function getAddressInputValue() {
  return addressInput.value.trim();
}

function validateAddressInput() {
  const address = getAddressInputValue();

  if (!address) {
    alert("Please enter an address first.");
    return null;
  }

  return address;
}

function formatHistoryItem(item) {
  return `
    <div style="padding: 12px 0; border-bottom: 1px solid #d8dfeb;">
      <p><strong>Action:</strong> ${item.action}</p>
      <p><strong>Amount:</strong> ${item.amount}</p>
      <p><strong>Status:</strong> ${item.status}</p>
      <p><strong>Tx Hash:</strong> ${item.tx_hash}</p>
      <p><strong>Created At:</strong> ${item.created_at}</p>
    </div>
  `;
}

function renderHistory(items) {
  if (!items || items.length === 0) {
    historyResult.innerHTML = `<p class="muted">No history found for this address.</p>`;
    return;
  }

  historyResult.innerHTML = items.map(formatHistoryItem).join("");
}

loadBalanceButton.addEventListener("click", async () => {
  const address = validateAddressInput();
  if (!address) {
    return;
  }

  try {
    console.log("Loading balance for:", address);

    const response = await getBalance(address);
    console.log("Balance response:", response);

    resultAddress.textContent = response.data?.address ?? address;
    resultBalance.textContent = response.data?.balance ?? "No balance";
  } catch (error) {
    console.error("Cannot load balance:", error);
    alert("Cannot load balance. Please check the address or backend.");
  }
});

loadHistoryButton.addEventListener("click", async () => {
  const address = validateAddressInput();
  if (!address) {
    return;
  }

  try {
    console.log("Loading history for:", address);

    const response = await getHistory(address);
    console.log("History response:", response);

    renderHistory(response.data?.items ?? []);
  } catch (error) {
    console.error("Cannot load history:", error);
    alert("Cannot load history. Please check the address or backend.");
  }
});

submitDepositButton.addEventListener("click", async () => {
  const amount = Number(depositAmount.value);

  if (!amount || amount <= 0) {
    alert("Please enter a valid deposit amount.");
    return;
  }

  try {
    console.log("Submitting deposit:", amount);
    const response = await deposit(amount);
    console.log("Deposit response:", response);

    alert(response.message ?? "Deposit submitted");
  } catch (error) {
    console.error("Cannot deposit:", error);
    alert("Cannot deposit. Please check the amount or backend.");
  }
});

submitWithdrawButton.addEventListener("click", async () => {
  const amount = Number(withdrawAmount.value);

  if (!amount || amount <= 0) {
    alert("Please enter a valid withdraw amount.");
    return;
  }

  try {
    console.log("Submitting withdraw:", amount);
    const response = await withdraw(amount);
    console.log("Withdraw response:", response);

    alert(response.message ?? "Withdraw submitted");
  } catch (error) {
    console.error("Cannot withdraw:", error);
    alert("Cannot withdraw. Please check the amount or backend.");
  }
});
