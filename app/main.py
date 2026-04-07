from fastapi import FastAPI
import yfinance as yf

app = FastAPI()

@app.get("/")
def root():
    return {"message": "Stock API OK 🚀"}

@app.get("/saham/{kode}")
def get_saham(kode: str):
    ticker = yf.Ticker(f"{kode}.JK")
    info = ticker.info

    return {
        "kode": kode.upper(),
        "nama": info.get("longName"),
        "harga": info.get("currentPrice"),
        "target": info.get("targetMeanPrice"),
        "recommendation": info.get("recommendationKey")
    }