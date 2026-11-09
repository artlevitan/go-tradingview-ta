// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
)

const defaultScannerURL = "https://scanner.tradingview.com/symbol"

var (
	defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}

	// ErrNilTradingView indicates a nil destination passed to Get or GetContext.
	ErrNilTradingView = errors.New("tradingview: nil TradingView receiver")
	// ErrNilContext indicates a nil context passed to GetContext.
	ErrNilContext = errors.New("tradingview: nil context")
	// ErrInvalidSymbol indicates a symbol outside the EXCHANGE:SYMBOL format.
	ErrInvalidSymbol = errors.New("tradingview: symbol must be in EXCHANGE:SYMBOL format")
	// ErrInvalidInterval indicates an unsupported nonempty GetContext interval.
	ErrInvalidInterval = errors.New("tradingview: unsupported interval")
	// ErrInvalidResponse indicates an incomplete or invalid scanner response.
	ErrInvalidResponse = errors.New("tradingview: invalid scanner response")

	// DefaultClient is used by the TradingView methods. Configure it before
	// use; do not change it while requests are in flight.
	DefaultClient = Client{HTTPClient: defaultHTTPClient, BaseURL: defaultScannerURL}
)

// HTTPError describes an unsuccessful scanner HTTP response. Use errors.As to
// inspect it. The client does not automatically retry failed requests.
type HTTPError struct {
	StatusCode int    // HTTP response status code.
	Body       string // Response body, truncated to at most 4 KiB plus an ellipsis.
	RetryAfter string // Unparsed Retry-After header, if provided by the server.
}

func (e *HTTPError) Error() string {
	message := fmt.Sprintf("tradingview: HTTP %d (%s)", e.StatusCode, http.StatusText(e.StatusCode))
	if e.Body != "" {
		message += ": " + e.Body
	}
	return message
}

// Client fetches technical-analysis snapshots from TradingView.
// Its zero value is ready to use. A Client can be shared by concurrent requests
// to separate destinations, provided its configuration is not changed during use.
type Client struct {
	// HTTPClient executes requests. If nil, a shared client with a 10-second
	// timeout is used. A custom client is used as configured, including its timeout.
	HTTPClient *http.Client
	// BaseURL overrides the scanner URL. It must be an absolute HTTP(S) URL.
	// Existing query parameters are preserved, except symbol and fields.
	BaseURL string
}

// Get populates ta with a snapshot for symbol at interval.
// For compatibility with v1.3, an empty or unsupported interval requests daily
// data. Use GetContext for strict interval validation and request cancellation.
// On error, ta is unchanged.
func (ta *TradingView) Get(symbol, interval string) error {
	return DefaultClient.Get(ta, symbol, interval)
}

// GetContext populates ta with a snapshot for symbol at interval using ctx.
// See Client.GetContext for validation and availability semantics.
func (ta *TradingView) GetContext(ctx context.Context, symbol, interval string) error {
	return DefaultClient.GetContext(ctx, ta, symbol, interval)
}

// Get populates ta using c. It is the compatibility wrapper for GetContext:
// unsupported intervals request daily data, and context.Background is used.
// On error, ta is unchanged.
func (c *Client) Get(ta *TradingView, symbol, interval string) error {
	if _, err := intervalSuffix(interval); err != nil {
		interval = Interval1Day
	}
	return c.GetContext(context.Background(), ta, symbol, interval)
}

// GetContext populates ta with a snapshot for symbol at interval using c.
// Symbol must have the form "EXCHANGE:SYMBOL", for example "BINANCE:BTCUSDT".
// An empty interval requests daily data; an unsupported nonempty interval
// returns ErrInvalidInterval. The context must not be nil.
//
// Null indicator values are preserved in Raw and excluded from recommendations
// using SignalUnavailable. Malformed or incomplete responses return
// ErrInvalidResponse. On any error, the previous snapshot is unchanged.
// The destination must not be read or written concurrently with this call.
func (c *Client) GetContext(ctx context.Context, ta *TradingView, symbol, interval string) error {
	if ctx == nil {
		return ErrNilContext
	}
	if ta == nil {
		return ErrNilTradingView
	}
	if err := validateSymbol(symbol); err != nil {
		return err
	}
	suffix, err := intervalSuffix(interval)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	req, err := c.newRequest(ctx, symbol, suffix)
	if err != nil {
		return err
	}
	res, err := c.httpClient().Do(req)
	if err != nil {
		return fmt.Errorf("tradingview: send request: %w", err)
	}
	response, err := decodeResponse(res, suffix)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if interval == "" {
		interval = Interval1Day
	}
	next := TradingView{
		Raw: &RawValues{values: response.raw},
		Metadata: Metadata{
			Symbol:     symbol,
			Interval:   interval,
			FetchedAt:  time.Now().UTC(),
			BarTime:    response.barTime,
			UpdateMode: response.updateMode,
		},
	}
	next.populate(response.values)
	*ta = next
	return nil
}

func validateSymbol(symbol string) error {
	exchange, ticker, ok := strings.Cut(symbol, ":")
	if !ok || exchange == "" || ticker == "" || strings.Contains(ticker, ":") ||
		strings.ContainsFunc(symbol, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) {
		return ErrInvalidSymbol
	}
	return nil
}

func intervalSuffix(interval string) (string, error) {
	switch interval {
	case "", Interval1Day:
		return "", nil
	case Interval1Min, Interval5Min, Interval15Min, Interval30Min,
		Interval1Hour, Interval2Hour, Interval4Hour, Interval1Week, Interval1Month:
		return "|" + interval, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrInvalidInterval, interval)
	}
}

func (c *Client) newRequest(ctx context.Context, symbol, suffix string) (*http.Request, error) {
	baseURL := defaultScannerURL
	if c != nil && c.BaseURL != "" {
		baseURL = c.BaseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("tradingview: parse scanner URL: %w", err)
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, errors.New("tradingview: scanner URL must be an absolute HTTP(S) URL")
	}
	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("tradingview: parse scanner query: %w", err)
	}
	fields := append(fieldsForInterval(suffix), "time"+suffix, "update_mode")
	params.Set("symbol", symbol)
	params.Set("fields", strings.Join(fields, ","))
	u.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("tradingview: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) httpClient() *http.Client {
	if c != nil && c.HTTPClient != nil {
		return c.HTTPClient
	}
	return defaultHTTPClient
}
