import { deposit } from "../api.js";
import { getNumberInputValue, showError, showSuccess } from "../utils.js";

function bindDepositHandler(elements) {
  elements.submitDepositButton.addEventListener("click", async () => {
    const amount = getNumberInputValue(elements.depositAmount);

    if (!amount || amount <= 0) {
      showError("Please enter a valid deposit amount.");
      return;
    }

    try {
      const response = await deposit(amount);
      showSuccess(response.message ?? "Deposit submitted");
    } catch (error) {
      console.error("Cannot deposit:", error);
      showError("Cannot deposit. Please check the amount or backend.");
    }
  });
}

export { bindDepositHandler };
