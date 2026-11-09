//go:build integration

// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import "testing"

// These tests make real public scanner requests and must be explicitly enabled:
// go test -tags=integration -run TestLive -count=1 ./...
func TestLiveTradingView(t *testing.T) {
	for _, interval := range []string{
		Interval1Min, Interval5Min, Interval15Min, Interval30Min,
		Interval1Hour, Interval2Hour, Interval4Hour, Interval1Day,
		Interval1Week, Interval1Month,
	} {
		t.Run(interval, func(t *testing.T) {
			assertLiveSnapshot(t, "BINANCE:BTCUSDT", interval)
		})
	}
	for _, symbol := range []string{"BINANCE:ETHUSDT", "NASDAQ:AAPL", "OANDA:EURUSD"} {
		t.Run(symbol, func(t *testing.T) {
			assertLiveSnapshot(t, symbol, Interval1Day)
		})
	}
}

func assertLiveSnapshot(t *testing.T, symbol, interval string) {
	t.Helper()
	var ta TradingView
	if err := ta.GetContext(t.Context(), symbol, interval); err != nil {
		t.Fatal(err)
	}
	if !ta.HasValue("close") || ta.Value.Prices.Close <= 0 || len(ta.Raw.Map()) != 83 {
		t.Fatalf("incomplete live snapshot: close=%v, fields=%d", ta.Value.Prices.Close, len(ta.Raw.Map()))
	}
	if ta.Metadata.BarTime.IsZero() || ta.Metadata.UpdateMode == "" {
		t.Fatalf("missing live metadata: %+v", ta.Metadata)
	}
	for _, tt := range []struct {
		field  string
		signal int
	}{
		{"EMA200", ta.Recommend.MovingAverages.EMA200},
		{"SMA200", ta.Recommend.MovingAverages.SMA200},
	} {
		if !ta.HasValue(tt.field) && tt.signal != SignalUnavailable {
			t.Errorf("unavailable %s generated signal %d", tt.field, tt.signal)
		}
	}
	t.Logf("%s %s: close=%v, global=%+v, update_mode=%s, bar=%s", symbol, interval,
		ta.Value.Prices.Close, ta.Recommend.Global, ta.Metadata.UpdateMode, ta.Metadata.BarTime)
}
