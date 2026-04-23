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

function renderHistory(items, container) {
  if (!items || items.length === 0) {
    container.innerHTML = `<p class="muted">No history found for this address.</p>`;
    return;
  }

  container.innerHTML = items.map(formatHistoryItem).join("");
}

export { renderHistory };