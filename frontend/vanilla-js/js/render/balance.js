function renderBalance(address, balance, elements) {
    elements.resultAddress.textContent = address ?? "-";
    elements.resultBalance.textContent = balance ?? "No balance";
}

export { renderBalance };