// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"os"
	"reflect"
	"strings"
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
	for _, spec := range recommendationSpecs {
		for _, field := range spec.fields {
			if field == "close" {
				continue
			} // Null close invalidates the whole snapshot.
			t.Run(spec.path+"/"+field, func(t *testing.T) {
				response := completeResponse("|60")
				response[field+"|60"] = nil
				client := fixtureClient(t, func(*http.Request) string { return encodeFixture(t, response) })
				var ta TradingView
				if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", Interval1Hour); err != nil {
					t.Fatal(err)
				}
				if ta.HasValue(field) {
					t.Errorf("null %s must be unavailable", field)
				}
				if got := int(snapshotField(t, ta.Recommend, spec.path).Int()); got != SignalUnavailable {
					t.Errorf("null %s must make %s unavailable, got %d", field, spec.path, got)
				}
				assertSnapshotContract(t, &ta)
			})
		}
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

// These tables describe the public field-to-scanner contract independently of
// the population code. Distinct input values catch swapped or omitted fields.
var valueSpecs = []struct{ path, field string }{
	{"Global.Summary", "Recommend.All"},
	{"Global.Oscillators", "Recommend.Other"},
	{"Global.MA", "Recommend.MA"},
	{"Oscillators.RSI", "RSI"},
	{"Oscillators.StochK", "Stoch.K"},
	{"Oscillators.CCI", "CCI20"},
	{"Oscillators.ADX.Value", "ADX"},
	{"Oscillators.ADX.PlusDI", "ADX+DI"},
	{"Oscillators.ADX.MinusDI", "ADX-DI"},
	{"Oscillators.ADX.PlusDI1", "ADX+DI[1]"},
	{"Oscillators.ADX.MinusDI1", "ADX-DI[1]"},
	{"Oscillators.AO.Value", "AO"},
	{"Oscillators.AO.Prev1", "AO[1]"},
	{"Oscillators.AO.Prev2", "AO[2]"},
	{"Oscillators.Mom", "Mom"},
	{"Oscillators.MACD.Macd", "MACD.macd"},
	{"Oscillators.MACD.Signal", "MACD.signal"},
	{"Oscillators.StochRSI", "Stoch.RSI.K"},
	{"Oscillators.WR", "W.R"},
	{"Oscillators.BBP", "BBPower"},
	{"Oscillators.UO", "UO"},
	{"MovingAverages.EMA10", "EMA10"},
	{"MovingAverages.SMA10", "SMA10"},
	{"MovingAverages.EMA20", "EMA20"},
	{"MovingAverages.SMA20", "SMA20"},
	{"MovingAverages.EMA30", "EMA30"},
	{"MovingAverages.SMA30", "SMA30"},
	{"MovingAverages.EMA50", "EMA50"},
	{"MovingAverages.SMA50", "SMA50"},
	{"MovingAverages.EMA100", "EMA100"},
	{"MovingAverages.SMA100", "SMA100"},
	{"MovingAverages.EMA200", "EMA200"},
	{"MovingAverages.SMA200", "SMA200"},
	{"MovingAverages.Ichimoku", "Ichimoku.BLine"},
	{"MovingAverages.VWMA", "VWMA"},
	{"MovingAverages.HullMA", "HullMA9"},
	{"Pivots.Classic.Middle", "Pivot.M.Classic.Middle"},
	{"Pivots.Classic.R1", "Pivot.M.Classic.R1"},
	{"Pivots.Classic.R2", "Pivot.M.Classic.R2"},
	{"Pivots.Classic.R3", "Pivot.M.Classic.R3"},
	{"Pivots.Classic.S1", "Pivot.M.Classic.S1"},
	{"Pivots.Classic.S2", "Pivot.M.Classic.S2"},
	{"Pivots.Classic.S3", "Pivot.M.Classic.S3"},
	{"Pivots.Fibonacci.Middle", "Pivot.M.Fibonacci.Middle"},
	{"Pivots.Fibonacci.R1", "Pivot.M.Fibonacci.R1"},
	{"Pivots.Fibonacci.R2", "Pivot.M.Fibonacci.R2"},
	{"Pivots.Fibonacci.R3", "Pivot.M.Fibonacci.R3"},
	{"Pivots.Fibonacci.S1", "Pivot.M.Fibonacci.S1"},
	{"Pivots.Fibonacci.S2", "Pivot.M.Fibonacci.S2"},
	{"Pivots.Fibonacci.S3", "Pivot.M.Fibonacci.S3"},
	{"Pivots.Camarilla.Middle", "Pivot.M.Camarilla.Middle"},
	{"Pivots.Camarilla.R1", "Pivot.M.Camarilla.R1"},
	{"Pivots.Camarilla.R2", "Pivot.M.Camarilla.R2"},
	{"Pivots.Camarilla.R3", "Pivot.M.Camarilla.R3"},
	{"Pivots.Camarilla.S1", "Pivot.M.Camarilla.S1"},
	{"Pivots.Camarilla.S2", "Pivot.M.Camarilla.S2"},
	{"Pivots.Camarilla.S3", "Pivot.M.Camarilla.S3"},
	{"Pivots.Woodie.Middle", "Pivot.M.Woodie.Middle"},
	{"Pivots.Woodie.R1", "Pivot.M.Woodie.R1"},
	{"Pivots.Woodie.R2", "Pivot.M.Woodie.R2"},
	{"Pivots.Woodie.R3", "Pivot.M.Woodie.R3"},
	{"Pivots.Woodie.S1", "Pivot.M.Woodie.S1"},
	{"Pivots.Woodie.S2", "Pivot.M.Woodie.S2"},
	{"Pivots.Woodie.S3", "Pivot.M.Woodie.S3"},
	{"Pivots.Demark.Middle", "Pivot.M.Demark.Middle"},
	{"Pivots.Demark.R1", "Pivot.M.Demark.R1"},
	{"Pivots.Demark.S1", "Pivot.M.Demark.S1"},
	{"Prices.Close", "close"},
	{"Prices.High", "high"},
	{"Prices.Low", "low"},
}

var recommendationSpecs = []struct {
	path   string
	fields []string
}{
	{"Global.Summary", []string{"Recommend.All"}},
	{"Global.Oscillators", []string{"Recommend.Other"}},
	{"Global.MA", []string{"Recommend.MA"}},
	{"Oscillators.RSI", []string{"RSI", "RSI[1]"}},
	{"Oscillators.StochK", []string{"Stoch.K", "Stoch.D", "Stoch.K[1]", "Stoch.D[1]"}},
	{"Oscillators.CCI", []string{"CCI20", "CCI20[1]"}},
	{"Oscillators.ADX", []string{"ADX", "ADX+DI", "ADX-DI", "ADX+DI[1]", "ADX-DI[1]"}},
	{"Oscillators.AO", []string{"AO", "AO[1]", "AO[2]"}},
	{"Oscillators.Mom", []string{"Mom", "Mom[1]"}},
	{"Oscillators.MACD", []string{"MACD.macd", "MACD.signal"}},
	{"Oscillators.StochRSI", []string{"Rec.Stoch.RSI", "Stoch.RSI.K"}},
	{"Oscillators.WR", []string{"Rec.WR", "W.R"}},
	{"Oscillators.BBP", []string{"Rec.BBPower", "BBPower"}},
	{"Oscillators.UO", []string{"Rec.UO", "UO"}},
	{"MovingAverages.EMA10", []string{"EMA10", "close"}},
	{"MovingAverages.SMA10", []string{"SMA10", "close"}},
	{"MovingAverages.EMA20", []string{"EMA20", "close"}},
	{"MovingAverages.SMA20", []string{"SMA20", "close"}},
	{"MovingAverages.EMA30", []string{"EMA30", "close"}},
	{"MovingAverages.SMA30", []string{"SMA30", "close"}},
	{"MovingAverages.EMA50", []string{"EMA50", "close"}},
	{"MovingAverages.SMA50", []string{"SMA50", "close"}},
	{"MovingAverages.EMA100", []string{"EMA100", "close"}},
	{"MovingAverages.SMA100", []string{"SMA100", "close"}},
	{"MovingAverages.EMA200", []string{"EMA200", "close"}},
	{"MovingAverages.SMA200", []string{"SMA200", "close"}},
	{"MovingAverages.Ichimoku", []string{"Rec.Ichimoku", "Ichimoku.BLine"}},
	{"MovingAverages.VWMA", []string{"Rec.VWMA", "VWMA"}},
	{"MovingAverages.HullMA", []string{"Rec.HullMA9", "HullMA9"}},
}

func snapshotField(t *testing.T, value any, path string) reflect.Value {
	t.Helper()
	field := reflect.ValueOf(value)
	for name := range strings.SplitSeq(path, ".") {
		field = field.FieldByName(name)
		if !field.IsValid() {
			t.Fatalf("public snapshot field %q does not exist", path)
		}
	}
	return field
}

func assertSnapshotContract(t *testing.T, ta *TradingView) {
	t.Helper()
	raw := ta.Raw.Map()
	if len(raw) != 83 {
		t.Fatalf("snapshot contains %d numeric fields, want 83", len(raw))
	}
	for field, value := range raw {
		if value != nil && (math.IsNaN(*value) || math.IsInf(*value, 0)) {
			t.Errorf("non-finite raw value for %s", field)
		}
	}
	for _, spec := range valueSpecs {
		want := 0.0
		if number, ok := ta.Raw.Lookup(spec.field); ok {
			want = number
		}
		if got := snapshotField(t, ta.Value, spec.path).Float(); got != want {
			t.Errorf("Value.%s = %g, raw %s = %g", spec.path, got, spec.field, want)
		}
	}
	for _, spec := range recommendationSpecs {
		available := true
		for _, field := range spec.fields {
			available = available && ta.HasValue(field)
		}
		got := int(snapshotField(t, ta.Recommend, spec.path).Int())
		if !available {
			if got != SignalUnavailable {
				t.Errorf("Recommend.%s = %d despite unavailable input", spec.path, got)
			}
			continue
		}
		bound := SignalBuy
		if strings.HasPrefix(spec.path, "Global.") {
			bound = SignalStrongBuy
		}
		if got < -bound || got > bound {
			t.Errorf("Recommend.%s = %d outside the available signal range", spec.path, got)
		}
	}
}

func TestSnapshotFieldMapping(t *testing.T) {
	response := completeResponse("|60")
	for i, field := range indicatorFields {
		switch {
		case strings.HasPrefix(field, "Recommend."):
			response[field+"|60"] = float64(i+1) / 10
		case strings.HasPrefix(field, "Rec."):
			response[field+"|60"] = float64(i%3 - 1)
		default:
			response[field+"|60"] = float64(i + 1)
		}
	}
	client := fixtureClient(t, func(*http.Request) string { return encodeFixture(t, response) })
	var ta TradingView
	if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", Interval1Hour); err != nil {
		t.Fatal(err)
	}
	assertSnapshotContract(t, &ta)
}

func TestEveryMissingScannerFieldIsRejected(t *testing.T) {
	for _, field := range indicatorFields {
		t.Run(field, func(t *testing.T) {
			response := completeResponse("|60")
			delete(response, field+"|60")
			client := fixtureClient(t, func(*http.Request) string { return encodeFixture(t, response) })
			ta := TradingView{Value: Values{Prices: PriceValues{Close: 123}}}
			if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", Interval1Hour); !errors.Is(err, ErrInvalidResponse) {
				t.Fatalf("missing %s must be rejected, got %v", field, err)
			}
			if ta.Value.Prices.Close != 123 || ta.Raw != nil {
				t.Error("incomplete response changed the previous snapshot")
			}
		})
	}
}

func TestSnapshotWithoutAnalysisHistory(t *testing.T) {
	response := completeResponse("|60")
	for _, field := range indicatorFields {
		response[field+"|60"] = nil
	}
	response["close|60"] = 100.0
	delete(response, "time|60")
	delete(response, "update_mode")
	client := fixtureClient(t, func(*http.Request) string { return encodeFixture(t, response) })
	var ta TradingView
	if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", Interval1Hour); err != nil {
		t.Fatal(err)
	}
	assertSnapshotContract(t, &ta)
	if !ta.Metadata.BarTime.IsZero() || ta.Metadata.UpdateMode != "" {
		t.Errorf("missing metadata must remain unknown: %+v", ta.Metadata)
	}
}
