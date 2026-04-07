# 🤖📊 WhatsApp Stock Alert System

Sistem monitoring & notifikasi saham Indonesia berbasis **WhatsApp Bot**, menggunakan arsitektur microservices dengan Docker.

Project ini menggabungkan:

- 📡 Monitoring & alerting
- 📊 Data saham (Yahoo Finance)
- 💬 Notifikasi real-time ke WhatsApp

---

## 🧠 Arsitektur

```text
Alertmanager → Webhook → Stock API → Yahoo Finance
                                ↓
                            WhatsApp
```

---

## 🧱 Services

### 🔔 Alertmanager

Mengelola rule alert dan trigger notifikasi.

---

### 🧠 Webhook

- Menerima alert dari Alertmanager
- Memproses logic (format message, parsing, dll)
- Mengambil data saham dari Stock API

---

### 📊 Stock API

Service terpisah berbasis FastAPI untuk:

- Ambil data saham IDX (Indonesia)
- Integrasi dengan Yahoo Finance (`yfinance`)

---

### 💬 WAHA (WhatsApp HTTP API)

- Mengirim pesan ke WhatsApp
- Menjadi gateway komunikasi bot

---

## 📁 Project Structure

```bash
.
├── alertmanager/
│   └── alertmanager.yml
│
├── webhook/
│   └── (logic handler)
│
├── stock-api/
│   ├── app/
│   │   └── main.py
│   ├── Dockerfile
│   └── requirements.txt
│
├── waha-data/
│
└── docker-compose.yml
```

---

## ⚙️ Setup & Run

### 1. Jalankan semua service

```bash
docker-compose up -d --build
```

---

### 2. Akses service

| Service        | URL                   |
| -------------- | --------------------- |
| Stock API      | http://localhost:3010 |
| WAHA Dashboard | http://localhost:3005 |
| Webhook        | http://localhost:3001 |
| Alertmanager   | http://localhost:9093 |

---

## 📡 API Usage (Stock API)

### Get data saham

```bash
GET /saham/{kode}
```

Contoh:

```bash
GET /saham/DEWA
```

Response:

```json
{
  "kode": "DEWA",
  "nama": "PT Darma Henwa Tbk",
  "harga": 476.0,
  "target": 931.24,
  "recommendation": "strong_buy"
}
```

---

## 💬 WhatsApp Integration

Webhook akan:

1. Menerima trigger (alert / command)
2. Ambil data dari Stock API
3. Kirim ke WhatsApp via WAHA

Contoh output di WhatsApp:

```
📊 PT Darma Henwa Tbk
Harga: 476
Target: 931
Rekomendasi: strong_buy
```

---

## 🔗 Internal Communication

Gunakan nama service Docker:

```bash
http://stock-api:8000/saham/DEWA
http://waha:3000/api/sendText
```

---

## ⚠️ Notes

- Saham Indonesia otomatis menggunakan suffix `.JK`
- Data berasal dari Yahoo Finance (tidak real-time tick)
- `yfinance` bisa lambat → disarankan caching

---

## 🚀 Features (Current & Planned)

### ✅ Current

- Ambil data saham Indonesia
- Kirim notifikasi ke WhatsApp
- Integrasi Alertmanager

---

### 🔜 Planned

- 🔔 Price alert (trigger harga tertentu)
- 📊 Chart saham (kirim gambar ke WA)
- ⏱️ Scheduled report (harian)
- 🤖 Command dari WhatsApp (`harga dewa`)
- ⚡ Redis caching
- 🧠 AI analisis saham

---

## 🧠 Insight

Project ini dirancang sebagai:

> **Modular Notification Platform**

Bisa dikembangkan untuk:

- Monitoring server
- Crypto alerts
- CI/CD notifications
- Business alerts

---

## 👨‍💻 Author

Akbar Rahmat M
