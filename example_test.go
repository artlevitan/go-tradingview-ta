// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	tradingview "github.com/artlevitan/go-tradingview-ta"
)

func ExampleTradingView_Get() {
	server := httptest.NewServer(http.HandlerFunc(exampleScanner))
	defer server.Close()

	oldClient := tradingview.DefaultClient
	tradingview.DefaultClient = tradingview.Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	}
	defer func() {
		tradingview.DefaultClient = oldClient
	}()

	var ta tradingview.TradingView
	if err := ta.Get("BINANCE:BTCUSDT", tradingview.Interval1Hour); err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(ta.Recommend.Global.Summary)
	fmt.Println(ta.Value.Prices.Close)
	// Output:
	// 1
	// 70998.71
}

func ExampleClient_Get() {
	server := httptest.NewServer(http.HandlerFunc(exampleScanner))
	defer server.Close()

	client := tradingview.Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	}

	var ta tradingview.TradingView
	if err := client.Get(&ta, "BINANCE:BTCUSDT", tradingview.Interval1Hour); err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(ta.Recommend.Global.Summary)
	fmt.Println(ta.Value.Prices.Close)
	// Output:
	// 1
	// 70998.71
}

func ExampleClient_GetContext() {
	server := httptest.NewServer(http.HandlerFunc(exampleScanner))
	defer server.Close()
	client := tradingview.Client{HTTPClient: server.Client(), BaseURL: server.URL}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var ta tradingview.TradingView
	if err := client.GetContext(ctx, &ta, "BINANCE:BTCUSDT", tradingview.Interval1Hour); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(ta.Value.Prices.Close)
	fmt.Println(ta.Metadata.UpdateMode)
	// Output:
	// 70998.71
	// streaming
}

func ExampleTradingView_HasValue() {
	var ta tradingview.TradingView
	if err := json.Unmarshal([]byte(`{"Raw":{"MACD.macd":0,"EMA200":null}}`), &ta); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(ta.HasValue("MACD.macd"))
	fmt.Println(ta.HasValue("EMA200"))
	// Output:
	// true
	// false
}

// Examples return all requested fields, with nulls for unavailable indicators.
func exampleScanner(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]any)
	for field := range strings.SplitSeq(r.URL.Query().Get("fields"), ",") {
		response[field] = nil
	}
	response["Recommend.All|60"] = 0.3
	response["Recommend.Other|60"] = 0.0
	response["Recommend.MA|60"] = 0.8
	response["close|60"] = 70998.71
	response["update_mode"] = "streaming"
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
