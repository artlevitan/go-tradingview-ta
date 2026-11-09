// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestTradingView_GetAllIntervals(t *testing.T) {
	ta := &TradingView{}
	symbol := "BINANCE:BTCUSDT"

	intervals := []struct {
		name     string
		interval string
		suffix   string
	}{
		{"Interval1min", Interval1Min, "|1"},
		{"Interval5min", Interval5Min, "|5"},
		{"Interval15min", Interval15Min, "|15"},
		{"Interval30min", Interval30Min, "|30"},
		{"Interval1hour", Interval1Hour, "|60"},
		{"Interval2hour", Interval2Hour, "|120"},
		{"Interval4hour", Interval4Hour, "|240"},
		{"Interval1day", Interval1Day, ""},
		{"Interval1week", Interval1Week, "|1W"},
		{"Interval1month", Interval1Month, "|1M"},
		{"default", "", ""},
	}

	for _, tt := range intervals {
		t.Run(tt.name, func(t *testing.T) {
			client := fixtureClient(t, func(r *http.Request) string {
				fields := strings.Split(r.URL.Query().Get("fields"), ",")
				if !slices.Contains(fields, "close"+tt.suffix) || !slices.Contains(fields, "time"+tt.suffix) || !slices.Contains(fields, "update_mode") {
					t.Errorf("incorrect interval or metadata fields: %v", fields)
				}
				if len(fields) != 85 {
					t.Errorf("got %d fields, want 83 indicators and 2 metadata fields", len(fields))
				}
				return encodeFixture(t, completeResponse(tt.suffix))
			})
			err := client.Get(ta, symbol, tt.interval)
			if err != nil {
				t.Fatalf("Interval %s: expected no error, got %v", tt.interval, err)
			}
			if ta.Value.Prices.Close <= 0 {
				t.Fatalf("Interval %s: expected positive close price, got %v", tt.interval, ta.Value.Prices.Close)
			}
			wantInterval := tt.interval
			if wantInterval == "" {
				wantInterval = Interval1Day
			}
			if ta.Metadata.Interval != wantInterval || ta.Metadata.Symbol != symbol || !ta.HasValue("close") {
				t.Fatalf("incorrect snapshot metadata: %+v", ta.Metadata)
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
			if !slices.Contains(fields, expected) {
				http.Error(w, "missing expected field "+expected, http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		response := completeResponse("|60")
		for field, value := range map[string]float64{
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
		} {
			response[field] = value
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
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
	err := client.Get(ta, "BINANCE:BTCUSDT", Interval1Hour)
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected HTTPError with status 502, got %v", err)
	}
}

func TestTradingView_InvalidInputDoesNotSendRequest(t *testing.T) {
	client := fixtureClient(t, func(*http.Request) string {
		t.Error("invalid input triggered an HTTP request")
		return `{}`
	})
	for _, symbol := range []string{"", ":", ":BTCUSDT", "BINANCE:", "BTCUSDT", "BINANCE:BTC:USDT", " BINANCE:BTCUSDT", "BINANCE:BTC USDT", "BINANCE:\u00a0BTCUSDT", "BINANCE:BTC\x00USDT"} {
		var ta TradingView
		if err := client.GetContext(context.Background(), &ta, symbol, Interval1Hour); !errors.Is(err, ErrInvalidSymbol) {
			t.Errorf("symbol=%q: expected ErrInvalidSymbol, got %v", symbol, err)
		}
	}
	for _, interval := range []string{"4h", "1d", "abc", "60 ", "2"} {
		var ta TradingView
		if err := client.GetContext(context.Background(), &ta, "BINANCE:BTCUSDT", interval); !errors.Is(err, ErrInvalidInterval) {
			t.Errorf("interval=%q: expected ErrInvalidInterval, got %v", interval, err)
		}
	}
}

func TestGetCompatibility(t *testing.T) {
	client := fixtureClient(t, func(r *http.Request) string {
		if !slices.Contains(strings.Split(r.URL.Query().Get("fields"), ","), "close") {
			t.Error("legacy Get must request daily data for an unknown interval")
		}
		return encodeFixture(t, completeResponse(""))
	})
	var ta TradingView
	if err := client.Get(&ta, "BINANCE:BTCUSDT", "4h"); err != nil {
		t.Fatal(err)
	}
	if ta.Metadata.Interval != Interval1Day {
		t.Errorf("legacy Get interval = %q, want daily", ta.Metadata.Interval)
	}
	if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", "4h"); !errors.Is(err, ErrInvalidInterval) {
		t.Errorf("GetContext must reject an unknown interval, got %v", err)
	}
}

func TestGetContextValidation(t *testing.T) {
	client := fixtureClient(t, func(*http.Request) string {
		t.Error("invalid or canceled call sent a request")
		return `{}`
	})
	var ta TradingView
	if err := client.GetContext(nil, &ta, "BINANCE:BTCUSDT", Interval1Hour); !errors.Is(err, ErrNilContext) {
		t.Errorf("nil context: got %v", err)
	}
	if err := client.GetContext(t.Context(), nil, "BINANCE:BTCUSDT", Interval1Hour); !errors.Is(err, ErrNilTradingView) {
		t.Errorf("nil destination: got %v", err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := client.GetContext(ctx, &ta, "BINANCE:BTCUSDT", Interval1Hour); !errors.Is(err, context.Canceled) {
		t.Errorf("canceled context: got %v", err)
	}
	ctx, cancel = context.WithDeadline(t.Context(), time.Unix(0, 0))
	defer cancel()
	if err := client.GetContext(ctx, &ta, "BINANCE:BTCUSDT", Interval1Hour); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expired deadline: got %v", err)
	}
}

func TestGetContextCancelsHTTP(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
	}))
	t.Cleanup(server.Close)
	client := Client{HTTPClient: server.Client(), BaseURL: server.URL}
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	ta := TradingView{Value: Values{Prices: PriceValues{Close: 123}}}
	result := make(chan error, 1)
	go func() { result <- client.GetContext(ctx, &ta, "BINANCE:BTCUSDT", Interval1Hour) }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("request did not reach the test server")
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) || ta.Value.Prices.Close != 123 {
			t.Fatalf("cancellation must preserve the snapshot: error=%v, ta=%+v", err, ta)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("request did not stop after cancellation")
	}
}

func TestGetContextCanceledResponsePreservesSnapshot(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	client := fixtureClient(t, func(*http.Request) string {
		cancel() // A custom transport may ignore cancellation and still return data.
		return encodeFixture(t, completeResponse("|60"))
	})
	ta := TradingView{Value: Values{Prices: PriceValues{Close: 123}}}
	if err := client.GetContext(ctx, &ta, "BINANCE:BTCUSDT", Interval1Hour); !errors.Is(err, context.Canceled) || ta.Value.Prices.Close != 123 {
		t.Fatalf("canceled response changed the snapshot: error=%v, ta=%+v", err, ta)
	}
}

func TestGetContextTransportError(t *testing.T) {
	want := errors.New("transport failed")
	client := Client{HTTPClient: &http.Client{Transport: testTransport(func(*http.Request) (*http.Response, error) {
		return nil, want
	})}}
	var ta TradingView
	if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", Interval1Hour); !errors.Is(err, want) {
		t.Errorf("transport error must remain wrapped, got %v", err)
	}
}

func TestRequestURL(t *testing.T) {
	for _, baseURL := range []string{
		"https://example.com/scanner",
		"https://example.com/scanner?",
		"https://example.com/scanner?market=crypto&",
		"https://example.com/scanner?market=crypto&symbol=OLD&fields=OLD#fragment",
	} {
		t.Run(baseURL, func(t *testing.T) {
			client := Client{BaseURL: baseURL}
			req, err := client.newRequest(t.Context(), "BINANCE:BTCUSDT", "|60")
			if err != nil {
				t.Fatal(err)
			}
			params := req.URL.Query()
			if req.Method != http.MethodGet || req.Header.Get("Accept") != "application/json" ||
				params.Get("symbol") != "BINANCE:BTCUSDT" || len(params["symbol"]) != 1 || len(params["fields"]) != 1 {
				t.Fatalf("invalid request: %+v", req)
			}
			if strings.Contains(baseURL, "market=crypto") && params.Get("market") != "crypto" {
				t.Fatal("existing query parameters were lost")
			}
		})
	}
	for _, baseURL := range []string{"://broken", "/relative", "ftp://example.com", "https://example.com?bad=%zz", "https://example.com/\x7f"} {
		client := fixtureClient(t, func(*http.Request) string {
			t.Error("invalid endpoint sent a request")
			return `{}`
		})
		client.BaseURL = baseURL
		var ta TradingView
		if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", Interval1Hour); err == nil {
			t.Errorf("invalid scanner URL %q was accepted", baseURL)
		}
	}
}

func TestDefaultClientMethods(t *testing.T) {
	previous := DefaultClient
	previousHTTP := defaultHTTPClient
	t.Cleanup(func() { DefaultClient, defaultHTTPClient = previous, previousHTTP })
	client := fixtureClient(t, func(*http.Request) string { return encodeFixture(t, completeResponse("|60")) })
	DefaultClient = *client
	defaultHTTPClient = client.HTTPClient
	var ta TradingView
	if err := ta.GetContext(t.Context(), "BINANCE:BTCUSDT", Interval1Hour); err != nil {
		t.Fatal(err)
	}
	if err := ta.Get("BINANCE:BTCUSDT", Interval1Hour); err != nil {
		t.Fatal(err)
	}
	// Both zero and nil clients use the default transport and endpoint.
	var zero Client
	var nilClient *Client
	for _, client := range []*Client{&zero, nilClient} {
		if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", Interval1Hour); err != nil {
			t.Fatal(err)
		}
	}
}

func TestClientConcurrentRequests(t *testing.T) {
	body := encodeFixture(t, completeResponse("|60"))
	client := fixtureClient(t, func(*http.Request) string { return body })
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			var ta TradingView
			if err := client.GetContext(t.Context(), &ta, "BINANCE:BTCUSDT", Interval1Hour); err != nil {
				t.Error(err)
			} else if !ta.HasValue("close") {
				t.Error("concurrent request lost its snapshot")
			}
		})
	}
	wg.Wait()
}

func TestDeadlineWhileReadingResponse(t *testing.T) {
	for _, mode := range []string{"context", "HTTP client"} {
		t.Run(mode, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.(http.Flusher).Flush()
				<-r.Context().Done() // Send headers, but never complete the body.
			}))
			t.Cleanup(server.Close)
			httpClient := server.Client()
			ctx := t.Context()
			if mode == "context" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 200*time.Millisecond)
				t.Cleanup(cancel)
			} else {
				httpClient.Timeout = 200 * time.Millisecond
			}
			client := Client{HTTPClient: httpClient, BaseURL: server.URL}
			ta := TradingView{Value: Values{Prices: PriceValues{Close: 123}}}
			err := client.GetContext(ctx, &ta, "BINANCE:BTCUSDT", Interval1Hour)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected the %s deadline to interrupt the body read, got %v", mode, err)
			}
			if ta.Value.Prices.Close != 123 || ta.Raw != nil {
				t.Fatal("timed-out response changed the previous snapshot")
			}
		})
	}
}
