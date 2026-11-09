// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	maxResponseBytes  = 1 << 20
	maxErrorBodyBytes = 4 << 10
)

var indicatorFields = [...]string{
	"Recommend.All",
	"Recommend.Other",
	"Recommend.MA",
	"RSI",
	"RSI[1]",
	"Stoch.K",
	"Stoch.D",
	"Stoch.K[1]",
	"Stoch.D[1]",
	"CCI20",
	"CCI20[1]",
	"ADX",
	"ADX+DI",
	"ADX-DI",
	"ADX+DI[1]",
	"ADX-DI[1]",
	"AO",
	"AO[1]",
	"AO[2]",
	"MACD.macd",
	"MACD.signal",
	"Mom",
	"Mom[1]",
	"Rec.Stoch.RSI",
	"Stoch.RSI.K",
	"Rec.WR",
	"W.R",
	"Rec.BBPower",
	"BBPower",
	"Rec.UO",
	"UO",
	"EMA10",
	"SMA10",
	"EMA20",
	"SMA20",
	"EMA30",
	"SMA30",
	"EMA50",
	"SMA50",
	"EMA100",
	"SMA100",
	"EMA200",
	"SMA200",
	"Rec.Ichimoku",
	"Ichimoku.BLine",
	"Rec.VWMA",
	"VWMA",
	"Rec.HullMA9",
	"HullMA9",
	"Pivot.M.Classic.S3",
	"Pivot.M.Classic.S2",
	"Pivot.M.Classic.S1",
	"Pivot.M.Classic.Middle",
	"Pivot.M.Classic.R1",
	"Pivot.M.Classic.R2",
	"Pivot.M.Classic.R3",
	"Pivot.M.Fibonacci.S3",
	"Pivot.M.Fibonacci.S2",
	"Pivot.M.Fibonacci.S1",
	"Pivot.M.Fibonacci.Middle",
	"Pivot.M.Fibonacci.R1",
	"Pivot.M.Fibonacci.R2",
	"Pivot.M.Fibonacci.R3",
	"Pivot.M.Camarilla.S3",
	"Pivot.M.Camarilla.S2",
	"Pivot.M.Camarilla.S1",
	"Pivot.M.Camarilla.Middle",
	"Pivot.M.Camarilla.R1",
	"Pivot.M.Camarilla.R2",
	"Pivot.M.Camarilla.R3",
	"Pivot.M.Woodie.S3",
	"Pivot.M.Woodie.S2",
	"Pivot.M.Woodie.S1",
	"Pivot.M.Woodie.Middle",
	"Pivot.M.Woodie.R1",
	"Pivot.M.Woodie.R2",
	"Pivot.M.Woodie.R3",
	"Pivot.M.Demark.S1",
	"Pivot.M.Demark.Middle",
	"Pivot.M.Demark.R1",
	"close",
	"high",
	"low",
}

func fieldsForInterval(suffix string) []string {
	fields := make([]string, len(indicatorFields))
	for i, field := range indicatorFields {
		fields[i] = field + suffix
	}
	return fields
}

type scannerResponse struct {
	values     indicatorValues
	raw        map[string]*float64
	barTime    time.Time
	updateMode string
}

func decodeResponse(res *http.Response, suffix string) (*scannerResponse, error) {
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(io.LimitReader(res.Body, maxErrorBodyBytes+1))
		if err != nil {
			return nil, fmt.Errorf("tradingview: read error response: %w", err)
		}
		truncated := len(body) > maxErrorBodyBytes
		if truncated {
			body = body[:maxErrorBodyBytes]
		}
		message := strings.TrimSpace(string(body))
		if truncated {
			message += "..."
		}
		return nil, &HTTPError{StatusCode: res.StatusCode, Body: message, RetryAfter: res.Header.Get("Retry-After")}
	}
	data, err := io.ReadAll(io.LimitReader(res.Body, maxResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("tradingview: read response: %w", err)
	}
	if len(data) > maxResponseBytes {
		return nil, fmt.Errorf("%w: response exceeds %d bytes", ErrInvalidResponse, maxResponseBytes)
	}
	// Go 1.27's JSON decoder rejects duplicate names and invalid UTF-8.
	var fields map[string]jsontext.Value
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, fmt.Errorf("%w: parse response: %w", ErrInvalidResponse, err)
	}
	response := &scannerResponse{
		values: make(indicatorValues, len(indicatorFields)),
		raw:    make(map[string]*float64, len(indicatorFields)),
	}
	for _, field := range indicatorFields {
		valueData, ok := fields[field+suffix]
		if !ok {
			return nil, fmt.Errorf("%w: missing field %q", ErrInvalidResponse, field+suffix)
		}
		var value *float64
		if err := json.Unmarshal(valueData, &value); err != nil {
			return nil, fmt.Errorf("%w: field %q: %w", ErrInvalidResponse, field+suffix, err)
		}
		response.raw[field] = value
		if value != nil {
			response.values[field] = *value
		}
	}
	if response.raw["close"] == nil {
		return nil, fmt.Errorf("%w: close is null", ErrInvalidResponse)
	}
	for _, field := range []string{"Recommend.All", "Recommend.Other", "Recommend.MA"} {
		if value := response.raw[field]; value != nil && (*value < -1 || *value > 1) {
			return nil, fmt.Errorf("%w: rating %q outside [-1, 1]", ErrInvalidResponse, field)
		}
	}
	for _, field := range []string{"Rec.Stoch.RSI", "Rec.WR", "Rec.BBPower", "Rec.UO", "Rec.Ichimoku", "Rec.VWMA", "Rec.HullMA9"} {
		if value := response.raw[field]; value != nil && *value != -1 && *value != 0 && *value != 1 {
			return nil, fmt.Errorf("%w: invalid signal %q", ErrInvalidResponse, field)
		}
	}
	if mode, ok := fields["update_mode"]; ok {
		if err := json.Unmarshal(mode, &response.updateMode); err != nil {
			return nil, fmt.Errorf("%w: update_mode: %w", ErrInvalidResponse, err)
		}
	}
	if barTime, ok := fields["time"+suffix]; ok {
		var timestamp *int64
		if err := json.Unmarshal(barTime, &timestamp); err != nil {
			return nil, fmt.Errorf("%w: time: %w", ErrInvalidResponse, err)
		}
		if timestamp != nil {
			response.barTime = time.Unix(*timestamp, 0).UTC()
			if year := response.barTime.Year(); year < 1 || year > 9999 {
				return nil, fmt.Errorf("%w: time outside the RFC 3339 year range", ErrInvalidResponse)
			}
		}
	}
	return response, nil
}
