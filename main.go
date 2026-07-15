package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type YahooResponse struct {
	QuoteSummary struct {
		Result []struct {
			FinancialData struct {
				CurrentPrice    PriceField `json:"currentPrice"`
				TargetMeanPrice PriceField `json:"targetMeanPrice"`
				RecommendKey    string     `json:"recommendationKey"`
			} `json:"financialData"`
			QuoteType struct {
				LongName string `json:"longName"`
			} `json:"quoteType"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"quoteSummary"`
}

type PriceField struct {
	Raw float64 `json:"raw"`
	Fmt string  `json:"fmt"`
}

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

var httpClient = &http.Client{Timeout: 10 * time.Second}

func fetchSaham(ctx context.Context, kode string) (*SahamResponse, error) {
	ticker := strings.ToUpper(kode) + ".JK"
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v10/finance/quoteSummary/%s?modules=financialData,quoteType",
		ticker,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var yr YahooResponse
	if err := json.Unmarshal(body, &yr); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	if yr.QuoteSummary.Error != nil {
		return nil, fmt.Errorf("yahoo error: %s", yr.QuoteSummary.Error.Description)
	}

	if len(yr.QuoteSummary.Result) == 0 {
		return &SahamResponse{
			Kode: strings.ToUpper(kode),
		}, nil
	}

	r := yr.QuoteSummary.Result[0]

	sr := &SahamResponse{
		Kode: strings.ToUpper(kode),
	}

	if r.QuoteType.LongName != "" {
		sr.Nama = r.QuoteType.LongName
	}
	if r.FinancialData.CurrentPrice.Raw != 0 {
		sr.Harga = r.FinancialData.CurrentPrice.Raw
	}
	if r.FinancialData.TargetMeanPrice.Raw != 0 {
		sr.Target = r.FinancialData.TargetMeanPrice.Raw
	}
	if r.FinancialData.RecommendKey != "" {
		sr.Recommendation = r.FinancialData.RecommendKey
	}

	return sr, nil
}

// fetchSahamBatch fetches multiple stocks concurrently using goroutines (like Promise.all in JS).
// Each stock is fetched in its own goroutine; results are collected via sync.WaitGroup.
// The parent context propagates cancellation — if the client disconnects, all in-flight fetches abort.
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

// handleSahamBatch handles GET /saham/batch?kode=BBCA,TLKM,ASII
// Spawns one goroutine per stock code, fetches all concurrently (like Promise.all).
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
