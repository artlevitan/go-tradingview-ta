// Copyright 2022-2024. All rights reserved.
// https://github.com/artlevitan/go-tradingview-ta
// v1.3.0

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
