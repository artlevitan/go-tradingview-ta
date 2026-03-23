// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	tradingview "github.com/artlevitan/go-tradingview-ta"
)

func ExampleTradingView_Get() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{
			"Recommend.All|60": 0.3,
			"Recommend.Other|60": 0,
			"Recommend.MA|60": 0.8,
			"close|60": 70998.71
		}`)
	}))
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{
			"Recommend.All|60": 0.3,
			"Recommend.Other|60": 0,
			"Recommend.MA|60": 0.8,
			"close|60": 70998.71
		}`)
	}))
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
