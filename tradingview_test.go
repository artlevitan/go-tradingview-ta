// Copyright 2022-2025. All rights reserved.
// https://github.com/artlevitan/go-tradingview-ta
// v1.3.4

package tradingview

import (
	"testing"
)

// Tests for all intervals
func TestTradingView_GetAllIntervals(t *testing.T) {
	// Create a TradingView object
	ta := &TradingView{}

	// Define the test symbol
	symbol := "BINANCE:BTCUSDT"

	// Table of test cases for all intervals
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

	// Run tests for each interval
	for _, tt := range intervals {
		t.Run(tt.name, func(t *testing.T) {
			err := ta.Get(symbol, tt.interval)
			if err != nil {
				t.Errorf("Interval %s: Expected no error, got %v", tt.interval, err)
			}
		})
	}
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

func Test_tvRsi(t *testing.T) {
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
			if got := tvRsi(tt.args.rsi, tt.args.rsi1); got != tt.want {
				t.Errorf("tvRsi() = %v, want %v", got, tt.want)
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

func Test_tvCci20(t *testing.T) {
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
			if got := tvCci20(tt.args.cci20, tt.args.cci201); got != tt.want {
				t.Errorf("tvCci20() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvAdx(t *testing.T) {
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
			if got := tvAdx(tt.args.adx, tt.args.adxpdi, tt.args.adxndi, tt.args.adxpdi1, tt.args.adxndi1); got != tt.want {
				t.Errorf("tvAdx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tvAo(t *testing.T) {
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
			if got := tvAo(tt.args.ao, tt.args.ao1, tt.args.ao2); got != tt.want {
				t.Errorf("tvAo() = %v, want %v", got, tt.want)
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

func Test_tvMacd(t *testing.T) {
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
			if got := tvMacd(tt.args.macd, tt.args.s); got != tt.want {
				t.Errorf("tvMacd() = %v, want %v", got, tt.want)
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

func Test_tvMa(t *testing.T) {
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
			if got := tvMa(tt.args.ma, tt.args.close); got != tt.want {
				t.Errorf("tvMa() = %v, want %v", got, tt.want)
			}
		})
	}
}
