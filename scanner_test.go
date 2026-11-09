// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error { b.closed = true; return nil }

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

func TestDecodeResponseClosesBody(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status int
		reader io.Reader
	}{
		{"success", http.StatusOK, strings.NewReader(encodeFixture(t, completeResponse("")))},
		{"invalid JSON", http.StatusOK, strings.NewReader("{")},
		{"HTTP error", http.StatusBadGateway, strings.NewReader("unavailable")},
		{"read error", http.StatusOK, failingReader{io.ErrUnexpectedEOF}},
		{"error body read failure", http.StatusBadGateway, failingReader{io.ErrUnexpectedEOF}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			body := &trackedBody{Reader: tt.reader}
			_, err := decodeResponse(&http.Response{StatusCode: tt.status, Body: body}, "")
			if !body.closed {
				t.Error("response body was not closed")
			}
			if strings.Contains(tt.name, "read") && !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Errorf("read error was not preserved: %v", err)
			}
		})
	}
}

func TestHTTPError(t *testing.T) {
	for _, tt := range []struct {
		name string
		body string
		want string
	}{
		{"body", "  rate limited\n", "rate limited"},
		{"empty", "", ""},
		{"bounded", strings.Repeat("x", maxErrorBodyBytes+100), strings.Repeat("x", maxErrorBodyBytes) + "..."},
	} {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{"Retry-After": []string{"120"}},
				Body:       io.NopCloser(strings.NewReader(tt.body)),
			}
			_, err := decodeResponse(res, "")
			var httpErr *HTTPError
			if !errors.As(err, &httpErr) {
				t.Fatalf("expected HTTPError, got %v", err)
			}
			if httpErr.StatusCode != http.StatusTooManyRequests || httpErr.RetryAfter != "120" || httpErr.Body != tt.want {
				t.Fatalf("incorrect HTTPError: %+v", httpErr)
			}
			if !strings.Contains(httpErr.Error(), "429") || !strings.Contains(httpErr.Error(), tt.want) {
				t.Errorf("error lacks status or details: %s", httpErr)
			}
		})
	}
}

func TestDecodeResponseSizeLimit(t *testing.T) {
	body := encodeFixture(t, completeResponse(""))
	for _, size := range []int{maxResponseBytes, maxResponseBytes + 1} {
		data := body + strings.Repeat(" ", size-len(body))
		_, err := decodeResponse(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(data))}, "")
		if size == maxResponseBytes && err != nil {
			t.Errorf("response at size limit was rejected: %v", err)
		}
		if size > maxResponseBytes && !errors.Is(err, ErrInvalidResponse) {
			t.Errorf("oversized response must fail, got %v", err)
		}
	}
}

func TestDecodeResponseStrictJSON(t *testing.T) {
	valid := encodeFixture(t, completeResponse("|60"))
	for _, tt := range []struct {
		name string
		body string
	}{
		{"duplicate key", strings.TrimSuffix(valid, "}") + `,"close|60":200}`},
		{"invalid UTF-8", strings.TrimSuffix(valid, "}") + ",\"unknown\":\"\xff\"}"},
		{"trailing object", valid + `{}`},
		{"overflow", strings.Replace(valid, `"close|60":100`, `"close|60":1e400`, 1)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeResponse(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(tt.body))}, "|60")
			if !errors.Is(err, ErrInvalidResponse) {
				t.Errorf("invalid JSON must fail, got %v", err)
			}
		})
	}
	// New server-side fields do not invalidate a complete known schema.
	body := strings.TrimSuffix(valid, "}") + `,"future_field":{"nested":"value"}}`
	if _, err := decodeResponse(&http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, "|60"); err != nil {
		t.Errorf("unknown extension field must be ignored: %v", err)
	}
}
