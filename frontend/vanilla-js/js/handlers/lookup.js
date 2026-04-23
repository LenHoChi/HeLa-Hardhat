import { getBalance, getHistory } from "../api.js";
import { renderBalance } from "../render/balance.js";
import { renderHistory } from "../render/history.js";
import { getAddressInputValue, showError } from "../utils.js";

function validateAddressInput(inputElement) {
  const address = getAddressInputValue(inputElement);

  if (!address) {
    showError("Please enter an address first.");
    return null;
  }

  return address;
}

function bindLookupHandlers(elements) {
  elements.loadBalanceButton.addEventListener("click", async () => {
    const address = validateAddressInput(elements.addressInput);
    if (!address) return;

    try {
      const response = await getBalance(address);

      renderBalance(
        response.data?.address ?? address,
        response.data?.balance ?? "No balance",
        {
          resultAddress: elements.resultAddress,
          resultBalance: elements.resultBalance,
        }
      );
    } catch (error) {
      console.error("Cannot load balance:", error);
      showError("Cannot load balance. Please check the address or backend.");
    }
  });

  elements.loadHistoryButton.addEventListener("click", async () => {
    const address = validateAddressInput(elements.addressInput);
    if (!address) return;

    try {
      const response = await getHistory(address);
      renderHistory(response.data?.items ?? [], elements.historyResult);
    } catch (error) {
      console.error("Cannot load history:", error);
      showError("Cannot load history. Please check the address or backend.");
    }
  });
}

export { bindLookupHandlers };
