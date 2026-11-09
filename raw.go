// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import json "encoding/json/v2"

// RawValues holds the numeric scanner fields for a snapshot, including nulls.
// Fields have no interval suffix. Use Lookup to read a number or Map to obtain
// an independent copy containing pointers for numbers and nil for nulls.
// Its zero value is empty. Read-only operations can be used concurrently.
// Do not unmarshal into a RawValues while other goroutines are reading it.
type RawValues struct {
	values map[string]*float64
}

// Lookup returns a field's number and reports whether it is available.
// Missing fields, null fields, and a nil receiver return (0, false).
func (raw *RawValues) Lookup(field string) (float64, bool) {
	if raw == nil || raw.values[field] == nil {
		return 0, false
	}
	return *raw.values[field], true
}

// Map returns a copy of all fields, including nulls. The returned map and its
// numeric pointers can be changed without affecting the snapshot.
// A nil or zero receiver returns nil.
func (raw *RawValues) Map() map[string]*float64 {
	if raw == nil || raw.values == nil {
		return nil
	}
	values := make(map[string]*float64, len(raw.values))
	for field, value := range raw.values {
		if value == nil {
			values[field] = nil
		} else {
			number := *value
			values[field] = &number
		}
	}
	return values
}

// MarshalJSON preserves null fields in the scanner's numeric JSON object.
func (raw RawValues) MarshalJSON() ([]byte, error) {
	if raw.values == nil {
		return []byte("null"), nil
	}
	return json.Marshal(raw.values)
}

// UnmarshalJSON decodes a numeric JSON object, retaining null fields.
// Invalid data leaves the existing values unchanged.
func (raw *RawValues) UnmarshalJSON(data []byte) error {
	var values map[string]*float64
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	raw.values = values
	return nil
}
