// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTradingView_GetAllIntervals(t *testing.T) {
	ta := &TradingView{}
	symbol := "BINANCE:BTCUSDT"

	intervals := []struct {
		name     string
		interval string
	}{
		{"Interval1min", Interval1Min},
		{"Interval5min", Interval5Min},
		{"Interval15min", Interval15Min},
		{"Interval30min", Interval30Min},
		{"Interval1hour", Interval1Hour},
		{"Interval2hour", Interval2Hour},
		{"Interval4hour", Interval4Hour},
		{"Interval1day", Interval1Day},
		{"Interval1week", Interval1Week},
		{"Interval1month", Interval1Month},
		{"default", ""}, // default
	}

	for _, tt := range intervals {
		t.Run(tt.name, func(t *testing.T) {
			err := ta.Get(symbol, tt.interval)
			if err != nil {
				t.Fatalf("Interval %s: expected no error, got %v", tt.interval, err)
			}
			if ta.Value.Prices.Close <= 0 {
				t.Fatalf("Interval %s: expected positive close price, got %v", tt.interval, ta.Value.Prices.Close)
			}
		})
	}
}

func TestTradingView_GetNilReceiver(t *testing.T) {
	var ta *TradingView

	err := ta.Get("BINANCE:BTCUSDT", Interval1Hour)
	if !errors.Is(err, ErrNilTradingView) {
		t.Fatalf("expected ErrNilTradingView, got %v", err)
	}
}

func TestTradingView_GetInvalidSymbol(t *testing.T) {
	ta := &TradingView{}

	err := ta.Get("BTCUSDT", Interval1Hour)
	if !errors.Is(err, ErrInvalidSymbol) {
		t.Fatalf("expected ErrInvalidSymbol, got %v", err)
	}
}

func TestTradingView_GetParsesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("symbol"); got != "BINANCE:BTCUSDT" {
			http.Error(w, "unexpected symbol", http.StatusBadRequest)
			return
		}

		fields := strings.Split(r.URL.Query().Get("fields"), ",")
		for _, expected := range []string{
			"Recommend.All|60",
			"Recommend.Other|60",
			"Recommend.MA|60",
			"ADX|60",
			"ADX+DI|60",
			"ADX-DI|60",
			"close|60",
		} {
			if !contains(fields, expected) {
				http.Error(w, "missing expected field "+expected, http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]float64{
			"Recommend.All|60":   0.35454545454545455,
			"Recommend.Other|60": -0.09090909090909091,
			"Recommend.MA|60":    0.8,
			"ADX|60":             34.379682334489544,
			"ADX+DI|60":          31.5847767930833,
			"ADX-DI|60":          19.08558943782413,
			"ADX+DI[1]|60":       18,
			"ADX-DI[1]|60":       21,
			"close|60":           70998.71,
			"high|60":            71050,
			"low|60":             70900,
		}); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	}

	ta := &TradingView{}
	if err := client.Get(ta, "BINANCE:BTCUSDT", Interval1Hour); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ta.Value.Global.Summary != 0.35454545454545455 {
		t.Fatalf("unexpected summary value: %v", ta.Value.Global.Summary)
	}
	if ta.Recommend.Global.Summary != SignalBuy {
		t.Fatalf("unexpected summary recommendation: %v", ta.Recommend.Global.Summary)
	}
	if ta.Recommend.Global.Oscillators != SignalNeutral {
		t.Fatalf("unexpected oscillators recommendation: %v", ta.Recommend.Global.Oscillators)
	}
	if ta.Recommend.Global.MA != SignalStrongBuy {
		t.Fatalf("unexpected MA recommendation: %v", ta.Recommend.Global.MA)
	}
	if ta.Value.Oscillators.ADX.MinusDI != 19.08558943782413 {
		t.Fatalf("unexpected ADX-DI value: %v", ta.Value.Oscillators.ADX.MinusDI)
	}
	if ta.Recommend.Oscillators.ADX != SignalBuy {
		t.Fatalf("unexpected ADX recommendation: %v", ta.Recommend.Oscillators.ADX)
	}
	if ta.Value.Prices.Close != 70998.71 {
		t.Fatalf("unexpected close price: %v", ta.Value.Prices.Close)
	}
}

func TestTradingView_GetUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"unavailable"}`, http.StatusBadGateway)
	}))
	defer server.Close()

	client := Client{
		HTTPClient: server.Client(),
		BaseURL:    server.URL,
	}

	ta := &TradingView{}
	if err := client.Get(ta, "BINANCE:BTCUSDT", Interval1Hour); err == nil {
		t.Fatal("expected error for non-200 response")
	}
}

func Test_key(t *testing.T) {
	tests := []struct {
		name         string
		indicator    string
		dataInterval string
		want         string
	}{
		{name: "with interval", indicator: "close%s", dataInterval: "|60", want: "close|60"},
		{name: "default interval", indicator: "close%s", dataInterval: "", want: "close"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := key(tt.indicator, tt.dataInterval); got != tt.want {
				t.Fatalf("key() = %q, want %q", got, tt.want)
			}
		})
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func Test_tvComputeRecommend(t *testing.T) {
	type args struct {
		v float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Strong Buy", args{0.6}, SignalStrongBuy},      // v between 0.5 and 1
		{"Test Buy", args{0.3}, SignalBuy},                   // v between 0.1 and 0.5
		{"Test Neutral Positive", args{0.1}, SignalNeutral},  // v between -0.1 and 0.1
		{"Test Neutral Negative", args{-0.1}, SignalNeutral}, // v between -0.1 and 0.1
		{"Test Strong Sell", args{-0.6}, SignalStrongSell},   // v between -1 and -0.5
		{"Test Sell", args{-0.3}, SignalSell},                // v between -0.5 and -0.1
		{"Test Neutral Edge", args{2}, SignalNeutral},        // value outside of range
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvComputeRecommend(tt.args.v); got != tt.want {
				t.Errorf("tvComputeRecommend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvRSI(t *testing.T) {
	type args struct {
		rsi  float64
		rsi1 float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Buy", args{29, 28}, SignalBuy},         // RSI < 30 and rsi1 < rsi
		{"Test Sell", args{71, 72}, SignalSell},       // RSI > 70 and rsi1 > rsi
		{"Test Neutral", args{50, 50}, SignalNeutral}, // Default case (neutral)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvRSI(tt.args.rsi, tt.args.rsi1); got != tt.want {
				t.Errorf("tvRSI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvStoch(t *testing.T) {
	type args struct {
		k  float64
		d  float64
		k1 float64
		d1 float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Buy", args{19, 18, 17, 18}, SignalBuy},         // k, d < 20, k > d, k1 < d1
		{"Test Sell", args{81, 82, 83, 82}, SignalSell},       // k, d > 80, k < d, k1 > d1
		{"Test Neutral", args{50, 50, 50, 50}, SignalNeutral}, // Default case (neutral)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvStoch(tt.args.k, tt.args.d, tt.args.k1, tt.args.d1); got != tt.want {
				t.Errorf("tvStoch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvCCI20(t *testing.T) {
	type args struct {
		cci20  float64
		cci201 float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Buy", args{-101, -102}, SignalBuy},   // CCI20 < -100 and CCI20 > CCI201
		{"Test Sell", args{101, 102}, SignalSell},   // CCI20 > 100 and CCI20 < CCI201
		{"Test Neutral", args{0, 0}, SignalNeutral}, // Default case (neutral)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvCCI20(tt.args.cci20, tt.args.cci201); got != tt.want {
				t.Errorf("tvCCI20() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvADX(t *testing.T) {
	type args struct {
		adx     float64
		adxpdi  float64
		adxndi  float64
		adxpdi1 float64
		adxndi1 float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Buy", args{21, 30, 20, 19, 21}, SignalBuy},         // ADX > 20, adxpdi1 < adxndi1, adxpdi > adxndi
		{"Test Sell", args{21, 20, 30, 31, 29}, SignalSell},       // ADX > 20, adxpdi1 > adxndi1, adxpdi < adxndi
		{"Test Neutral", args{19, 20, 20, 20, 20}, SignalNeutral}, // ADX <= 20
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvADX(tt.args.adx, tt.args.adxpdi, tt.args.adxndi, tt.args.adxpdi1, tt.args.adxndi1); got != tt.want {
				t.Errorf("tvADX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvAO(t *testing.T) {
	type args struct {
		ao  float64
		ao1 float64
		ao2 float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Buy", args{1, -1, 0}, SignalBuy},        // AO > 0, AO1 < 0
		{"Test Sell", args{-1, 1, 0}, SignalSell},      // AO < 0, AO1 > 0
		{"Test Neutral", args{0, 0, 0}, SignalNeutral}, // Default case (neutral)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvAO(tt.args.ao, tt.args.ao1, tt.args.ao2); got != tt.want {
				t.Errorf("tvAO() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvMom(t *testing.T) {
	type args struct {
		mom  float64
		mom1 float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Buy", args{2, 1}, SignalBuy},         // mom > mom1
		{"Test Sell", args{1, 2}, SignalSell},       // mom < mom1
		{"Test Neutral", args{1, 1}, SignalNeutral}, // Default case (neutral)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvMom(tt.args.mom, tt.args.mom1); got != tt.want {
				t.Errorf("tvMom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvMACD(t *testing.T) {
	type args struct {
		macd float64
		s    float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Buy", args{1, 0}, SignalBuy},         // MACD > Signal
		{"Test Sell", args{0, 1}, SignalSell},       // MACD < Signal
		{"Test Neutral", args{1, 1}, SignalNeutral}, // Default case (neutral)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvMACD(tt.args.macd, tt.args.s); got != tt.want {
				t.Errorf("tvMACD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvSimple(t *testing.T) {
	type args struct {
		v float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Buy", args{1}, SignalBuy},         // v == 1
		{"Test Sell", args{-1}, SignalSell},      // v == -1
		{"Test Neutral", args{0}, SignalNeutral}, // Default case (neutral)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvSimple(tt.args.v); got != tt.want {
				t.Errorf("tvSimple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvMA(t *testing.T) {
	type args struct {
		ma    float64
		close float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Test Buy", args{100, 101}, SignalBuy},         // ma < close
		{"Test Sell", args{101, 100}, SignalSell},       // ma > close
		{"Test Neutral", args{100, 100}, SignalNeutral}, // Default case (neutral)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tvMA(tt.args.ma, tt.args.close); got != tt.want {
				t.Errorf("tvMA() = %v, want %v", got, tt.want)
			}
		})
	}
}
