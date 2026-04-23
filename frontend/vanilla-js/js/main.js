import { API_BASE_URL, getBalance, getHistory } from "./api.js";

const addressInput = document.querySelector("#address-input");
const loadBalanceButton = document.querySelector("#load-balance-btn");
const loadHistoryButton = document.querySelector("#load-history-btn");
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
