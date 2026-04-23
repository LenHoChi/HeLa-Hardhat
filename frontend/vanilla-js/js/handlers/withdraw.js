import { withdraw } from "../api.js";
import { getNumberInputValue, showError, showSuccess } from "../utils.js";

function bindWithdrawHandler(elements) {
  elements.submitWithdrawButton.addEventListener("click", async () => {
    const amount = getNumberInputValue(elements.withdrawAmount);

    if (!amount || amount <= 0) {
      showError("Please enter a valid withdraw amount.");
      return;
    }

    try {
      const response = await withdraw(amount);
      showSuccess(response.message ?? "Withdraw submitted");
    } catch (error) {
      console.error("Cannot withdraw:", error);
      showError("Cannot withdraw. Please check the amount or backend.");
    }
  });
}

export { bindWithdrawHandler };
