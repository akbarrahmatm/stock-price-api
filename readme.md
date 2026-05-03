# 📈 Stock Price API

A lightweight REST API for fetching real-time Indonesian stock market (IDX) data, built with **FastAPI** and **yfinance**.

---

## 🚀 Features

- Real-time IDX stock price lookup
- Returns company name, current price, analyst target price, and recommendation
- Fast and minimal — powered by FastAPI

---

## 🛠️ Tech Stack

| Library                                            | Purpose                      |
| -------------------------------------------------- | ---------------------------- |
| [FastAPI](https://fastapi.tiangolo.com/)           | Web framework                |
| [yfinance](https://github.com/ranaroussi/yfinance) | Stock data via Yahoo Finance |
| [Uvicorn](https://www.uvicorn.org/)                | ASGI server                  |

---

## ⚙️ Installation

```bash
# Clone the repository
git clone https://github.com/akbarrahmatm/stock-price-api.git
cd stock-price-app

# Install dependencies
pip install -r requirements.txt

# Start the server
uvicorn app.main:app --reload
```

Server runs at: `http://localhost:8000`

---

## 📦 Requirements

```
fastapi
uvicorn
yfinance
```

> Save the above as `requirements.txt` in the root of your project.

---

## 📡 Endpoints

### `GET /`

Health check — verifies the API is up and running.

**Response:**

```json
{
  "message": "Stock API OK 🚀"
}
```

---

### `GET /saham/{kode}`

Fetch stock data by IDX ticker symbol.

**Path Parameter:**

| Parameter | Type     | Description                      |
| --------- | -------- | -------------------------------- |
| `kode`    | `string` | IDX stock ticker (without `.JK`) |

**Example Request:**

```
GET /saham/BBCA
```

**Example Response:**

```json
{
  "kode": "BBCA",
  "nama": "Bank Central Asia Tbk",
  "harga": 9350,
  "target": 10500,
  "recommendation": "buy"
}
```

**Popular IDX Tickers:**

| Code   | Company               |
| ------ | --------------------- |
| `BBCA` | Bank Central Asia     |
| `TLKM` | Telkom Indonesia      |
| `GOTO` | GoTo Gojek Tokopedia  |
| `ASII` | Astra International   |
| `BBRI` | Bank Rakyat Indonesia |

---

## 📂 Project Structure

```
└── stock-price-app
    └── app
        └── main.py
    └── requirements.txt
```

---

## 📝 Notes

- The `.JK` suffix (Jakarta Stock Exchange) is automatically appended to every ticker when querying Yahoo Finance.
- Data is sourced from **Yahoo Finance** via `yfinance` — a few minutes of delay may occur.
- `target` and `recommendation` are based on analyst consensus and may be `null` if unavailable.

---

## 📄 License

MIT License — free to use and modify.
