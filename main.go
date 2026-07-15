package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/AmpyFin/yfinance-go"
)

type SahamResponse struct {
	Kode           string      `json:"kode"`
	Nama           interface{} `json:"nama"`
	Harga          interface{} `json:"harga"`
	Target         interface{} `json:"target"`
	Recommendation interface{} `json:"recommendation"`
}

type sahamResult struct {
	Data  *SahamResponse `json:"data,omitempty"`
	Error string         `json:"error,omitempty"`
}

var yfClient = yfinance.NewClient()

func fetchSaham(ctx context.Context, kode string) (*SahamResponse, error) {
	ticker := strings.ToUpper(kode) + ".JK"
	sr := &SahamResponse{Kode: strings.ToUpper(kode)}

	quote, err := yfClient.FetchQuote(ctx, ticker, "stock-api")
	if err != nil {
		return nil, fmt.Errorf("quote: %w", err)
	}

	if quote.Security.Symbol != "" {
		sr.Nama = quote.Security.Symbol
	}
	if quote.RegularMarketPrice != nil {
		sr.Harga = float64(quote.RegularMarketPrice.Scaled) / float64(quote.RegularMarketPrice.Scale)
	}

	info, err := yfClient.FetchCompanyInfo(ctx, ticker, "stock-api")
	if err == nil && info.LongName != "" {
		sr.Nama = info.LongName
	}

	insights, err := yfClient.ScrapeAnalystInsights(ctx, ticker, "stock-api")
	if err == nil && insights != nil {
		for _, line := range insights.Lines {
			if line.Key == "target_price_mean" {
				if line.Value != nil {
					sr.Target = float64(line.Value.Scaled) / float64(line.Value.Scale)
				}
			}
			if line.Key == "recommendation_score" {
				if line.Value != nil {
					sr.Recommendation = float64(line.Value.Scaled) / float64(line.Value.Scale)
				}
			}
		}
	}

	return sr, nil
}

func fetchSahamBatch(ctx context.Context, kodeList []string) map[string]sahamResult {
	results := make(map[string]sahamResult, len(kodeList))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, kode := range kodeList {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			data, err := fetchSaham(ctx, k)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				results[strings.ToUpper(k)] = sahamResult{Error: err.Error()}
			} else {
				results[strings.ToUpper(k)] = sahamResult{Data: data}
			}
		}(kode)
	}

	wg.Wait()
	return results
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Stock API OK"})
}

func handleSaham(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/saham/"), "/")
	kode := parts[0]
	if kode == "" {
		http.Error(w, `{"error":"kode saham required"}`, http.StatusBadRequest)
		return
	}

	data, err := fetchSaham(r.Context(), kode)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func handleSahamBatch(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("kode")
	if raw == "" {
		http.Error(w, `{"error":"query param 'kode' required, e.g. ?kode=BBCA,TLKM"}`, http.StatusBadRequest)
		return
	}

	kodeList := strings.Split(raw, ",")
	for i := range kodeList {
		kodeList[i] = strings.TrimSpace(kodeList[i])
	}

	results := fetchSahamBatch(r.Context(), kodeList)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/saham/batch", handleSahamBatch)
	mux.HandleFunc("/saham/", handleSaham)
	return mux
}

func main() {
	mux := newMux()
	addr := ":8000"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
