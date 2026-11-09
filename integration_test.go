//go:build integration

// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

// These tests make real public scanner requests and must be explicitly enabled:
// go test -tags=integration -run TestLive -count=1 ./...
func TestLiveTradingView(t *testing.T) {
	intervals := []string{
		Interval1Min, Interval5Min, Interval15Min, Interval30Min,
		Interval1Hour, Interval2Hour, Interval4Hour, Interval1Day,
		Interval1Week, Interval1Month,
	}
	symbols := []string{
		"BINANCE:BTCUSDT", "COINBASE:BTCUSD", "KRAKEN:ETHUSD",
		"NASDAQ:AAPL", "NYSE:IBM", "XETR:SAP",
		"OANDA:EURUSD", "FX_IDC:USDJPY",
	}
	// Keep live checks sequential and modestly paced to avoid request bursts.
	pace := time.NewTicker(250 * time.Millisecond)
	defer pace.Stop()
	for _, symbol := range symbols {
		for _, interval := range intervals {
			t.Run(symbol+"/"+interval, func(t *testing.T) {
				select {
				case <-pace.C:
				case <-t.Context().Done():
					t.Fatal(t.Context().Err())
				}
				assertLiveSnapshot(t, symbol, interval)
			})
		}
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
	if ta.Metadata.Symbol != symbol || ta.Metadata.Interval != interval || ta.Metadata.FetchedAt.IsZero() {
		t.Fatalf("snapshot does not describe the requested instrument and interval: %+v", ta.Metadata)
	}
	assertSnapshotContract(t, &ta)
	t.Logf("%s %s: close=%v, global=%+v, update_mode=%s, bar=%s", symbol, interval,
		ta.Value.Prices.Close, ta.Recommend.Global, ta.Metadata.UpdateMode, ta.Metadata.BarTime)
}

func TestLiveNonexistentSymbol(t *testing.T) {
	ta := TradingView{Value: Values{Prices: PriceValues{Close: 123}}}
	err := ta.GetContext(t.Context(), "BINANCE:__GO_TRADINGVIEW_TA_MISSING__", Interval1Day)
	var statusErr *HTTPError
	if !errors.Is(err, ErrInvalidResponse) &&
		!(errors.As(err, &statusErr) && (statusErr.StatusCode == http.StatusNotFound || statusErr.StatusCode == http.StatusBadRequest)) {
		t.Fatalf("expected a real scanner rejection for the missing symbol, got %v", err)
	}
	if ta.Value.Prices.Close != 123 || ta.Raw != nil {
		t.Fatal("rejected symbol replaced the previous snapshot")
	}
	t.Logf("scanner rejected the nonexistent symbol: %v", err)
}
