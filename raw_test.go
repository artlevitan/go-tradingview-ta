// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"testing"
)

func TestRawValues(t *testing.T) {
	var raw RawValues
	if err := json.Unmarshal([]byte(`{"MACD.macd":0,"EMA10":12,"EMA200":null}`), &raw); err != nil {
		t.Fatal(err)
	}
	if value, ok := raw.Lookup("MACD.macd"); !ok || value != 0 {
		t.Error("a real zero must be available")
	}
	for _, field := range []string{"EMA200", "missing"} {
		if _, ok := raw.Lookup(field); ok {
			t.Errorf("%q should be unavailable", field)
		}
	}
	copy := raw.Map()
	*copy["EMA10"] = -1
	copy["EMA200"] = copy["EMA10"]
	delete(copy, "MACD.macd")
	if value, ok := raw.Lookup("EMA10"); !ok || value != 12 {
		t.Error("mutating a copied number changed the snapshot")
	}
	if value, ok := raw.Lookup("MACD.macd"); !ok || value != 0 {
		t.Error("mutating a copied map changed the snapshot")
	}
	if _, ok := raw.Lookup("EMA200"); ok {
		t.Error("mutating a copied null changed the snapshot")
	}
	if err := raw.UnmarshalJSON([]byte(`{"EMA10":"invalid"}`)); err == nil {
		t.Fatal("invalid JSON value accepted")
	}
	if value, _ := raw.Lookup("EMA10"); value != 12 {
		t.Error("failed unmarshal changed the original values")
	}
	for _, marshal := range []func(any) ([]byte, error){json.Marshal, func(v any) ([]byte, error) { return jsonv2.Marshal(v) }} {
		for _, value := range []any{raw, &raw} {
			data, err := marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			var decoded map[string]*float64
			if err := json.Unmarshal(data, &decoded); err != nil || len(decoded) != 3 || decoded["EMA200"] != nil || *decoded["EMA10"] != 12 {
				t.Fatalf("raw JSON did not round-trip: %s, error=%v", data, err)
			}
		}
	}
}

func TestRawValuesZero(t *testing.T) {
	var raw RawValues
	var nilRaw *RawValues
	for _, value := range []*RawValues{&raw, nilRaw} {
		if number, ok := value.Lookup("EMA10"); ok || number != 0 || value.Map() != nil {
			t.Error("empty RawValues must have no available fields")
		}
		data, err := json.Marshal(value)
		if err != nil || string(data) != "null" {
			t.Fatalf("empty RawValues JSON = %s, error=%v", data, err)
		}
	}
}

func TestTradingViewRemainsComparable(t *testing.T) {
	// v1.3 allowed comparisons and map keys; v1.4 must still compile these uses.
	var first, second TradingView
	if first != second {
		t.Error("zero snapshots must be equal")
	}
	cache := map[TradingView]bool{first: true}
	if !cache[second] {
		t.Error("TradingView must remain usable as a map key")
	}
}
