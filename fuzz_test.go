// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"testing"
)

func FuzzGetContext(f *testing.F) {
	monthly, err := os.ReadFile("testdata/btcusdt-1m.json")
	if err != nil {
		f.Fatal(err)
	}
	for _, seed := range [][]byte{monthly, []byte(`{}`), []byte(`null`), []byte(`{"close|1M":100}`)} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, body []byte) {
		client := fixtureClient(t, func(*http.Request) string { return string(body) })
		ta := TradingView{Value: Values{Prices: PriceValues{Close: 123}}}
		if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", Interval1Month); err != nil {
			if ta.Value.Prices.Close != 123 || ta.Raw != nil || !ta.Metadata.FetchedAt.IsZero() {
				t.Fatal("failed request mutated its destination")
			}
			return
		}
		values := ta.Raw.Map()
		if !ta.HasValue("close") || len(values) != 83 {
			t.Fatal("successful response has an incomplete schema")
		}
		for _, value := range values {
			if value != nil && (math.IsNaN(*value) || math.IsInf(*value, 0)) {
				t.Fatal("successful response contains a non-finite number")
			}
		}
		if _, err := json.Marshal(ta); err != nil {
			t.Fatalf("successful snapshot is not serializable: %v", err)
		}
	})
}
