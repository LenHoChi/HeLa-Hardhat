const API_BASE_URL = "http://localhost:8080";

export async function getBalance(address) {
  const response = await fetch(`${API_BASE_URL}/balance/${address}`);

  if (!response.ok) {
    throw new Error(`Request failed with status ${response.status}`);
  }

  return response.json();
}

export async function getHistory(address) {
  const response = await fetch(`${API_BASE_URL}/history/${address}`);

  if (!response.ok) {
    throw new Error(`Request failed with status ${response.status}`);
  }

  return response.json();
}

export async function deposit(amount) {
  const response = await fetch(`${API_BASE_URL}/deposit`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ amount }),
  });

  if (!response.ok) {
    throw new Error(`Request failed with status ${response.status}`);
  }

  return response.json();
}

export async function withdraw(amount) {
  const response = await fetch(`${API_BASE_URL}/withdraw`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ amount }),
  });

  if (!response.ok) {
    throw new Error(`Request failed with status ${response.status}`);
  }

  return response.json();
}

export { API_BASE_URL };
