# Stock Price API

A lightweight REST API for fetching real-time Indonesian stock market (IDX) data, built with **Go** and the **Yahoo Finance** API.

Compiled to a single static binary. Docker image is ~15 MB.

---

## Features

- Real-time IDX stock price lookup
- Batch endpoint — fetch multiple stocks concurrently via goroutines
- Context-aware HTTP calls with automatic cancellation
- Minimal Docker image using multi-stage build (`scratch`)

---

## Tech Stack

| Component | Purpose |
| --- | --- |
| Go `net/http` (stdlib) | HTTP server & router |
| Yahoo Finance API | Stock data source |
| `sync.WaitGroup` / goroutines | Concurrent batch fetching |

---

## Installation

### Run locally

```bash
git clone https://github.com/akbarrahmatm/akbarrahmatm-stock-api.git
cd akbarrahmatm-stock-api
go run .
```

Server runs at: `http://localhost:8000`

### Run with Docker

```bash
docker compose up --build
```

Server runs at: `http://localhost:3010`

---

## Endpoints

### `GET /`

Health check.

```json
{
  "message": "Stock API OK"
}
```

---

### `GET /saham/{kode}`

Fetch stock data by IDX ticker symbol.

| Parameter | Type | Description |
| --- | --- | --- |
| `kode` | `string` | IDX stock ticker (without `.JK`) |

**Example:** `GET /saham/BBCA`

```json
{
  "kode": "BBCA",
  "nama": "Bank Central Asia Tbk",
  "harga": 9350
}
```

---

### `GET /saham/batch?kode=BBCA,TLKM,ASII`

Fetch multiple stocks concurrently (goroutines, like `Promise.all` in JS).

**Example:** `GET /saham/batch?kode=BBCA,TLKM`

```json
{
  "BBCA": {
    "data": {
      "kode": "BBCA",
      "nama": "Bank Central Asia Tbk",
      "harga": 9350
    }
  },
  "TLKM": {
    "data": {
      "kode": "TLKM",
      "nama": "Telkom Indonesia Tbk",
      "harga": 3200
    }
  }
}
```

---

## Popular IDX Tickers

| Code | Company |
| --- | --- |
| `BBCA` | Bank Central Asia |
| `TLKM` | Telkom Indonesia |
| `GOTO` | GoTo Gojek Tokopedia |
| `ASII` | Astra International |
| `BBRI` | Bank Rakyat Indonesia |

---

## Testing

```bash
go test -v ./...
```

---

## Project Structure

```
.
├── main.go           # Server, handlers, Yahoo Finance client
├── main_test.go      # Tests (9 cases)
├── go.mod
├── Dockerfile        # Multi-stage build → scratch
└── docker-compose.yml
```

---

## Notes

- The `.JK` suffix (Jakarta Stock Exchange) is automatically appended to every ticker.
- Data is sourced from Yahoo Finance's unofficial API — a few minutes of delay may occur.

---

## License

MIT License — free to use and modify.
