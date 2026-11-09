// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type testTransport func(*http.Request) (*http.Response, error)

func (f testTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func fixtureClient(t *testing.T, body func(*http.Request) string) *Client {
	t.Helper()
	return &Client{HTTPClient: &http.Client{Transport: testTransport(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body(r)))}, nil
	})}}
}

func completeResponse(suffix string) map[string]any {
	response := make(map[string]any)
	for _, field := range fieldsForInterval(suffix) {
		response[field] = float64(0)
	}
	response["close"+suffix] = float64(100)
	response["time"+suffix] = int64(1794225600)
	response["update_mode"] = "streaming"
	return response
}

func encodeFixture(t *testing.T, response map[string]any) string {
	t.Helper()
	b, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
