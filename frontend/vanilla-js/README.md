# Hela Bank Vanilla JS UI

Simple frontend UI for the Hela Bank backend, built with plain HTML, CSS, and JavaScript.

This version is used to learn the basic frontend flow before moving to React.

## Features

- Load wallet balance
- Load transaction history
- Submit deposit
- Submit withdraw
- Call backend APIs with `fetch`

## Project Structure

```text
frontend/vanilla-js/
├── index.html
├── css/
│   └── style.css
├── js/
│   ├── api.js
│   ├── config.js
│   ├── main.js
│   ├── utils.js
│   ├── handlers/
│   └── render/
├── package.json
└── README.md
```

## Requirements

- Node.js
- Yarn
- Backend API running on `http://localhost:8080`

## Run

Install dependencies:

```bash
yarn install
```

Start the UI:

```bash
yarn dev
```

Open the URL printed by Vite, usually:

```text
http://localhost:5173
```

## Backend URL

The backend API URL is configured in:

```text
js/config.js
```

Default value:

```js
const API_BASE_URL = "http://localhost:8080";
```

## Notes

- This is a learning UI, so the code is intentionally simple.
- The same API logic can be reused later when building the React version.
- If the browser blocks requests, make sure CORS is enabled in the Go backend.
