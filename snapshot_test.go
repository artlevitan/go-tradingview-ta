// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestTradingView_MonthlySnapshot(t *testing.T) {
	// Recorded public scanner response, 2026-11-09. The official monthly
	// technicals page shows dashes for EMA200, SMA200, and Classic S3.
	b, err := os.ReadFile("testdata/btcusdt-1m.json")
	if err != nil {
		t.Fatal(err)
	}
	client := fixtureClient(t, func(*http.Request) string { return string(b) })
	var ta TradingView
	if err := client.Get(&ta, "BINANCE:BTCUSDT", Interval1Month); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"EMA200", "SMA200", "Pivot.M.Classic.S3"} {
		value, ok := ta.Raw.Map()[field]
		if !ok || value != nil || ta.HasValue(field) {
			t.Errorf("%s must remain unavailable: %v, present=%v", field, value, ok)
		}
	}
	if ta.Recommend.MovingAverages.EMA200 != SignalUnavailable || ta.Recommend.MovingAverages.SMA200 != SignalUnavailable {
		t.Fatalf("missing monthly averages generated signals: %+v", ta.Recommend.MovingAverages)
	}
	if got := ta.Recommend.Global; got != (GlobalRecommendations{Summary: SignalNeutral, Oscillators: SignalNeutral, MA: SignalSell}) {
		t.Errorf("monthly global ratings do not match TradingView: %+v", got)
	}
	// These are the counts displayed by the official technicals page.
	for _, tt := range []struct {
		name    string
		signals any
		want    map[int]int
	}{
		{"oscillators", ta.Recommend.Oscillators, map[int]int{SignalSell: 1, SignalNeutral: 9, SignalBuy: 1}},
		{"moving averages", ta.Recommend.MovingAverages, map[int]int{SignalSell: 7, SignalNeutral: 1, SignalBuy: 5, SignalUnavailable: 2}},
	} {
		counts := make(map[int]int)
		values := reflect.ValueOf(tt.signals)
		for _, field := range values.Fields() {
			counts[int(field.Int())]++
		}
		if !reflect.DeepEqual(counts, tt.want) {
			t.Errorf("%s counts: got %v, want %v", tt.name, counts, tt.want)
		}
	}
	b, err = json.Marshal(ta)
	if err != nil {
		t.Fatal(err)
	}
	var encoded struct{ Raw map[string]*float64 }
	if err := json.Unmarshal(b, &encoded); err != nil || encoded.Raw["EMA200"] != nil {
		t.Fatalf("JSON must preserve null EMA200: %s, error=%v", b, err)
	}
}

func TestTradingView_NullInputs(t *testing.T) {
	for _, tt := range []struct {
		field  string
		signal func(*TradingView) int
	}{
		{"Recommend.All", func(ta *TradingView) int { return ta.Recommend.Global.Summary }},
		{"RSI[1]", func(ta *TradingView) int { return ta.Recommend.Oscillators.RSI }},
		{"Stoch.D", func(ta *TradingView) int { return ta.Recommend.Oscillators.StochK }},
		{"CCI20[1]", func(ta *TradingView) int { return ta.Recommend.Oscillators.CCI }},
		{"ADX-DI[1]", func(ta *TradingView) int { return ta.Recommend.Oscillators.ADX }},
		{"AO[2]", func(ta *TradingView) int { return ta.Recommend.Oscillators.AO }},
		{"Mom[1]", func(ta *TradingView) int { return ta.Recommend.Oscillators.Mom }},
		{"MACD.signal", func(ta *TradingView) int { return ta.Recommend.Oscillators.MACD }},
		{"Rec.Stoch.RSI", func(ta *TradingView) int { return ta.Recommend.Oscillators.StochRSI }},
		{"Stoch.RSI.K", func(ta *TradingView) int { return ta.Recommend.Oscillators.StochRSI }},
		{"W.R", func(ta *TradingView) int { return ta.Recommend.Oscillators.WR }},
		{"BBPower", func(ta *TradingView) int { return ta.Recommend.Oscillators.BBP }},
		{"UO", func(ta *TradingView) int { return ta.Recommend.Oscillators.UO }},
		{"EMA10", func(ta *TradingView) int { return ta.Recommend.MovingAverages.EMA10 }},
		{"SMA200", func(ta *TradingView) int { return ta.Recommend.MovingAverages.SMA200 }},
		{"Ichimoku.BLine", func(ta *TradingView) int { return ta.Recommend.MovingAverages.Ichimoku }},
		{"VWMA", func(ta *TradingView) int { return ta.Recommend.MovingAverages.VWMA }},
		{"HullMA9", func(ta *TradingView) int { return ta.Recommend.MovingAverages.HullMA }},
	} {
		t.Run(tt.field, func(t *testing.T) {
			response := completeResponse("|60")
			response[tt.field+"|60"] = nil
			client := fixtureClient(t, func(*http.Request) string { return encodeFixture(t, response) })
			var ta TradingView
			if err := client.Get(&ta, "BINANCE:BTCUSDT", Interval1Hour); err != nil {
				t.Fatal(err)
			}
			if ta.HasValue(tt.field) || tt.signal(&ta) != SignalUnavailable {
				t.Errorf("null %s must produce SignalUnavailable, got %d", tt.field, tt.signal(&ta))
			}
			if !ta.HasValue("MACD.macd") || ta.Value.Oscillators.MACD.Macd != 0 {
				t.Error("a real zero must remain a valid numeric value")
			}
		})
	}
}

func TestTradingView_InvalidResponsesPreserveSnapshot(t *testing.T) {
	valid := completeResponse("|60")
	for _, tt := range []struct {
		name string
		body string
	}{
		{"empty object", `{}`},
		{"null", `null`},
		{"partial", `{"close|60":100,"EMA10|60":null}`},
		{"malformed", `{`},
		{"array", `[]`},
		{"wrong interval", encodeFixture(t, completeResponse(""))},
		{"missing symbol", `{"code":"symbol_not_exists","errmsg":"empty response"}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			body := encodeFixture(t, valid)
			client := fixtureClient(t, func(*http.Request) string { return body })
			var ta TradingView
			if err := client.Get(&ta, "BINANCE:BTCUSDT", Interval1Hour); err != nil {
				t.Fatal(err)
			}
			before := ta
			body = tt.body
			if err := client.Get(&ta, "BINANCE:BTCUSDT", Interval1Hour); !errors.Is(err, ErrInvalidResponse) {
				t.Errorf("expected ErrInvalidResponse, got %v", err)
			}
			if !reflect.DeepEqual(ta, before) {
				t.Error("an invalid response changed the previous snapshot")
			}
		})
	}
	for _, tt := range []struct {
		field string
		value any
	}{
		{"close|60", nil},
		{"RSI|60", "unknown"},
		{"EMA200|60", true},
		{"Recommend.All|60", 1.1},
		{"Recommend.Other|60", -1.1},
		{"Recommend.MA|60", 5},
		{"Rec.WR|60", 0.5},
		{"time|60", "not a timestamp"},
		{"time|60", int64(253402300800)}, // First second of year 10000 cannot be serialized as RFC 3339.
		{"update_mode", 123},
	} {
		t.Run(tt.field, func(t *testing.T) {
			response := completeResponse("|60")
			response[tt.field] = tt.value
			client := fixtureClient(t, func(*http.Request) string { return encodeFixture(t, response) })
			var ta TradingView
			if err := client.Get(&ta, "BINANCE:BTCUSDT", Interval1Hour); !errors.Is(err, ErrInvalidResponse) {
				t.Errorf("expected ErrInvalidResponse, got %v", err)
			}
		})
	}
}

func TestTradingView_MetadataAndSnapshotReplacement(t *testing.T) {
	response := completeResponse("|60")
	response["update_mode"] = "delayed_streaming_900"
	response["EMA200|60"] = 123.0
	client := fixtureClient(t, func(*http.Request) string { return encodeFixture(t, response) })
	var ta TradingView
	started := time.Now().UTC()
	if err := client.Get(&ta, "NASDAQ:AAPL", Interval1Hour); err != nil {
		t.Fatal(err)
	}
	if ta.Metadata.UpdateMode != "delayed_streaming_900" || !ta.Metadata.BarTime.Equal(time.Unix(1794225600, 0)) ||
		ta.Metadata.FetchedAt.Before(started) || ta.Metadata.FetchedAt.After(time.Now().UTC()) {
		t.Fatalf("incorrect freshness metadata: %+v", ta.Metadata)
	}
	previousRaw := ta.Raw
	response["EMA200|60"] = nil
	response["time|60"] = nil
	response["update_mode"] = nil
	if err := client.Get(&ta, "BINANCE:BTCUSDT", Interval1Hour); err != nil {
		t.Fatal(err)
	}
	previousValue, previousOK := previousRaw.Lookup("EMA200")
	if ta.HasValue("EMA200") || ta.Value.MovingAverages.EMA200 != 0 || ta.Recommend.MovingAverages.EMA200 != SignalUnavailable ||
		!ta.Metadata.BarTime.IsZero() || ta.Metadata.UpdateMode != "" || !previousOK || previousValue != 123 {
		t.Fatalf("new snapshot retained stale values: %+v", ta)
	}
	var empty TradingView
	var nilTA *TradingView
	if empty.HasValue("EMA200") || nilTA.HasValue("EMA200") || ta.HasValue("does_not_exist") {
		t.Error("missing fields must not report availability")
	}
}
