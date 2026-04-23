import { API_BASE_URL } from "./config.js";
import { bindDepositHandler } from "./handlers/deposit.js";
import { bindLookupHandlers } from "./handlers/lookup.js";
import { bindWithdrawHandler } from "./handlers/withdraw.js";
import { renderApiBaseUrl } from "./render/status.js";

const elements = {
  addressInput: document.querySelector("#address-input"),
  loadBalanceButton: document.querySelector("#load-balance-btn"),
  loadHistoryButton: document.querySelector("#load-history-btn"),
  submitDepositButton: document.querySelector("#submit-deposit-btn"),
  submitWithdrawButton: document.querySelector("#submit-withdraw-btn"),
  depositAmount: document.querySelector("#deposit-amount"),
  withdrawAmount: document.querySelector("#withdraw-amount"),
  resultAddress: document.querySelector("#result-address"),
  resultBalance: document.querySelector("#result-balance"),
  historyResult: document.querySelector("#history-result"),
  apiBaseUrl: document.querySelector("#api-base-url"),
};

renderApiBaseUrl(API_BASE_URL, elements.apiBaseUrl)

bindLookupHandlers(elements);
bindDepositHandler(elements);
bindWithdrawHandler(elements);
