// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultScannerURL = "https://scanner.tradingview.com/symbol"

var (
	defaultHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	// ErrNilTradingView reports that a nil *TradingView receiver was passed to Get.
	ErrNilTradingView = errors.New("tradingview: nil TradingView receiver")
	// ErrInvalidSymbol reports that a symbol does not have the form EXCHANGE:SYMBOL.
	ErrInvalidSymbol = errors.New("tradingview: symbol must be in EXCHANGE:SYMBOL format")
	// DefaultClient is used by (*TradingView).Get.
	DefaultClient = Client{
		HTTPClient: defaultHTTPClient,
		BaseURL:    defaultScannerURL,
	}
)

// Supported interval values for Get.
const (
	// Interval1Min requests 1-minute data.
	Interval1Min = "1"
	// Interval5Min requests 5-minute data.
	Interval5Min = "5"
	// Interval15Min requests 15-minute data.
	Interval15Min = "15"
	// Interval30Min requests 30-minute data.
	Interval30Min = "30"
	// Interval1Hour requests 1-hour data.
	Interval1Hour = "60"
	// Interval2Hour requests 2-hour data.
	Interval2Hour = "120"
	// Interval4Hour requests 4-hour data.
	Interval4Hour = "240"
	// Interval1Day requests 1-day data.
	Interval1Day = "1D"
	// Interval1Week requests 1-week data.
	Interval1Week = "1W"
	// Interval1Month requests 1-month data.
	Interval1Month = "1M"
)

// Normalized recommendation signals returned in TradingView.Recommend.
const (
	// SignalStrongBuy indicates a strong buy recommendation.
	SignalStrongBuy = 2
	// SignalBuy indicates a buy recommendation.
	SignalBuy = 1
	// SignalNeutral indicates a neutral recommendation.
	SignalNeutral = 0
	// SignalSell indicates a sell recommendation.
	SignalSell = -1
	// SignalStrongSell indicates a strong sell recommendation.
	SignalStrongSell = -2
)

// TradingView holds normalized recommendations and raw values returned by
// TradingView's scanner endpoint.
//
// Its zero value is ready to use.
type TradingView struct {
	// Recommend contains normalized recommendation signals.
	Recommend Recommendations
	// Value contains raw numeric indicator and price values.
	Value Values
}

// Recommendations groups normalized BUY, SELL, and NEUTRAL signals.
type Recommendations struct {
	// Global contains the combined recommendation groups.
	Global GlobalRecommendations
	// Oscillators contains oscillator recommendations.
	Oscillators OscillatorRecommendations
	// MovingAverages contains moving-average recommendations.
	MovingAverages MovingAverageRecommendations
}

// GlobalRecommendations stores the combined recommendation groups.
type GlobalRecommendations struct {
	Summary     int // Summary recommendation
	Oscillators int // Oscillators recommendation
	MA          int // Moving Averages recommendation
}

// OscillatorRecommendations stores normalized signals for oscillator indicators.
type OscillatorRecommendations struct {
	RSI      int // Relative Strength Index (14)
	StochK   int // Stochastic %K (14, 3, 3)
	CCI      int // Commodity Channel Index (20)
	ADX      int // Average Directional Index (14)
	AO       int // Awesome Oscillator
	Mom      int // Momentum (10)
	MACD     int // MACD Level (12, 26)
	StochRSI int // Stochastic RSI Fast (3, 3, 14, 14)
	WR       int // Williams Percent Range (14)
	BBP      int // Bull Bear Power
	UO       int // Ultimate Oscillator (7, 14, 28)
}

// MovingAverageRecommendations stores normalized signals for moving-average indicators.
type MovingAverageRecommendations struct {
	EMA10    int // Exponential Moving Average (EMA10)
	SMA10    int // Simple Moving Average (SMA10)
	EMA20    int // Exponential Moving Average (EMA20)
	SMA20    int // Simple Moving Average (SMA20)
	EMA30    int // Exponential Moving Average (EMA30)
	SMA30    int // Simple Moving Average (SMA30)
	EMA50    int // Exponential Moving Average (EMA50)
	SMA50    int // Simple Moving Average (SMA50)
	EMA100   int // Exponential Moving Average (EMA100)
	SMA100   int // Simple Moving Average (SMA100)
	EMA200   int // Exponential Moving Average (EMA200)
	SMA200   int // Simple Moving Average (SMA200)
	Ichimoku int // Ichimoku Base Line (9, 26, 52, 26)
	VWMA     int // Volume Weighted Moving Average (20)
	HullMA   int // Hull Moving Average (HullMA9)
}

// Values groups raw numeric values returned by TradingView.
type Values struct {
	// Global contains the raw combined recommendation values.
	Global GlobalValues
	// Oscillators contains raw oscillator values.
	Oscillators OscillatorValues
	// MovingAverages contains raw moving-average values.
	MovingAverages MovingAverageValues
	// Pivots contains pivot levels.
	Pivots PivotValues
	// Prices contains raw price values.
	Prices PriceValues
}

// GlobalValues stores the raw combined recommendation values.
type GlobalValues struct {
	Summary     float64 // Summary recommendation
	Oscillators float64 // Oscillators recommendation
	MA          float64 // Moving Averages recommendation
}

// OscillatorValues stores raw oscillator values.
type OscillatorValues struct {
	RSI      float64    // Relative Strength Index (14)
	StochK   float64    // Stochastic %K (14, 3, 3)
	CCI      float64    // Commodity Channel Index (20)
	ADX      ADXValues  // Average Directional Index (14)
	AO       AOValues   // Awesome Oscillator
	Mom      float64    // Momentum (10)
	MACD     MACDValues // MACD Level (12, 26)
	StochRSI float64    // Stochastic RSI Fast (3, 3, 14, 14)
	WR       float64    // Williams Percent Range (14)
	BBP      float64    // Bull Bear Power
	UO       float64    // Ultimate Oscillator (7, 14, 28)
}

// ADXValues stores Average Directional Index values.
type ADXValues struct {
	Value    float64 // ADX Value
	PlusDI   float64 // ADX+DI
	MinusDI  float64 // ADX-DI
	PlusDI1  float64 // ADX+DI[1]
	MinusDI1 float64 // ADX-DI[1]
}

// AOValues stores Awesome Oscillator values.
type AOValues struct {
	Value float64 // AO current value
	Prev1 float64 // AO[1]
	Prev2 float64 // AO[2]
}

// MACDValues stores MACD line values.
type MACDValues struct {
	Macd   float64 // MACD line
	Signal float64 // Signal line
}

// MovingAverageValues stores raw moving-average values.
type MovingAverageValues struct {
	EMA10    float64 // Exponential Moving Average (EMA10)
	SMA10    float64 // Simple Moving Average (SMA10)
	EMA20    float64 // Exponential Moving Average (EMA20)
	SMA20    float64 // Simple Moving Average (SMA20)
	EMA30    float64 // Exponential Moving Average (EMA30)
	SMA30    float64 // Simple Moving Average (SMA30)
	EMA50    float64 // Exponential Moving Average (EMA50)
	SMA50    float64 // Simple Moving Average (SMA50)
	EMA100   float64 // Exponential Moving Average (EMA100)
	SMA100   float64 // Simple Moving Average (SMA100)
	EMA200   float64 // Exponential Moving Average (EMA200)
	SMA200   float64 // Simple Moving Average (SMA200)
	Ichimoku float64 // Ichimoku Base Line (9, 26, 52, 26)
	VWMA     float64 // Volume Weighted Moving Average (20)
	HullMA   float64 // Hull Moving Average (HullMA9)
}

// PivotValues stores pivot levels for the supported pivot systems.
type PivotValues struct {
	Classic   ClassicPivotLevels
	Fibonacci FibonacciPivotLevels
	Camarilla CamarillaPivotLevels
	Woodie    WoodiePivotLevels
	Demark    DemarkPivotLevels
}

// ClassicPivotLevels stores Classic pivot levels.
type ClassicPivotLevels struct {
	Middle float64 // Classic Pivot Middle (Pivot.M.Classic.Middle)
	R1     float64 // Resistance 1 (Pivot.M.Classic.R1)
	R2     float64 // Resistance 2 (Pivot.M.Classic.R2)
	R3     float64 // Resistance 3 (Pivot.M.Classic.R3)
	S1     float64 // Support 1 (Pivot.M.Classic.S1)
	S2     float64 // Support 2 (Pivot.M.Classic.S2)
	S3     float64 // Support 3 (Pivot.M.Classic.S3)
}

// FibonacciPivotLevels stores Fibonacci pivot levels.
type FibonacciPivotLevels struct {
	Middle float64 // Fibonacci Pivot Middle (Pivot.M.Fibonacci.Middle)
	R1     float64 // Resistance 1 (Pivot.M.Fibonacci.R1)
	R2     float64 // Resistance 2 (Pivot.M.Fibonacci.R2)
	R3     float64 // Resistance 3 (Pivot.M.Fibonacci.R3)
	S1     float64 // Support 1 (Pivot.M.Fibonacci.S1)
	S2     float64 // Support 2 (Pivot.M.Fibonacci.S2)
	S3     float64 // Support 3 (Pivot.M.Fibonacci.S3)
}

// CamarillaPivotLevels stores Camarilla pivot levels.
type CamarillaPivotLevels struct {
	Middle float64 // Camarilla Pivot Middle (Pivot.M.Camarilla.Middle)
	R1     float64 // Resistance 1 (Pivot.M.Camarilla.R1)
	R2     float64 // Resistance 2 (Pivot.M.Camarilla.R2)
	R3     float64 // Resistance 3 (Pivot.M.Camarilla.R3)
	S1     float64 // Support 1 (Pivot.M.Camarilla.S1)
	S2     float64 // Support 2 (Pivot.M.Camarilla.S2)
	S3     float64 // Support 3 (Pivot.M.Camarilla.S3)
}

// WoodiePivotLevels stores Woodie pivot levels.
type WoodiePivotLevels struct {
	Middle float64 // Woodie Pivot Middle (Pivot.M.Woodie.Middle)
	R1     float64 // Resistance 1 (Pivot.M.Woodie.R1)
	R2     float64 // Resistance 2 (Pivot.M.Woodie.R2)
	R3     float64 // Resistance 3 (Pivot.M.Woodie.R3)
	S1     float64 // Support 1 (Pivot.M.Woodie.S1)
	S2     float64 // Support 2 (Pivot.M.Woodie.S2)
	S3     float64 // Support 3 (Pivot.M.Woodie.S3)
}

// DemarkPivotLevels stores Demark pivot levels.
type DemarkPivotLevels struct {
	Middle float64 // Demark Pivot Middle (Pivot.M.Demark.Middle)
	R1     float64 // Resistance 1 (Pivot.M.Demark.R1)
	S1     float64 // Support 1 (Pivot.M.Demark.S1)
}

// PriceValues stores the close, high, and low values returned by TradingView.
type PriceValues struct {
	Close float64 // Closing price
	High  float64 // Highest price
	Low   float64 // Lowest price
}

// Client fetches technical-analysis data from the TradingView scanner endpoint.
//
// If HTTPClient is nil, Get uses a default client with a 10 second timeout.
// If BaseURL is empty, Get uses TradingView's public scanner endpoint.
// The zero value of Client is ready to use.
type Client struct {
	// HTTPClient is used to make requests.
	HTTPClient *http.Client
	// BaseURL is the scanner endpoint base URL.
	BaseURL string
}

// Get populates ta with recommendations and raw indicator values for symbol
// at interval.
//
// Symbol must have the form "EXCHANGE:SYMBOL", for example "BINANCE:BTCUSDT".
// An empty or unknown interval is treated as daily data.
func (ta *TradingView) Get(symbol, interval string) error {
	return DefaultClient.Get(ta, symbol, interval)
}

// Get populates ta with recommendations and raw indicator values for symbol
// at interval using c.
//
// Symbol must have the form "EXCHANGE:SYMBOL", for example "BINANCE:BTCUSDT".
// An empty or unknown interval is treated as daily data.
func (c *Client) Get(ta *TradingView, symbol, interval string) error {
	if ta == nil {
		return ErrNilTradingView
	}
	if err := validateSymbol(symbol); err != nil {
		return err
	}

	dataInterval := intervalSuffix(interval)
	responseMap, err := c.getResponseMap(symbol, dataInterval)
	if err != nil {
		return err
	}

	ta.populate(responseMap, dataInterval)
	return nil
}

func validateSymbol(symbol string) error {
	if strings.Count(symbol, ":") != 1 {
		return ErrInvalidSymbol
	}
	return nil
}

func intervalSuffix(interval string) string {
	switch interval {
	case Interval1Min:
		return "|1"
	case Interval5Min:
		return "|5"
	case Interval15Min:
		return "|15"
	case Interval30Min:
		return "|30"
	case Interval1Hour:
		return "|60"
	case Interval2Hour:
		return "|120"
	case Interval4Hour:
		return "|240"
	case Interval1Week:
		return "|1W"
	case Interval1Month:
		return "|1M"
	default:
		return ""
	}
}

func fieldsForInterval(dataInterval string) []string {
	return []string{
		fmt.Sprintf("Recommend.All%s", dataInterval),
		fmt.Sprintf("Recommend.Other%s", dataInterval),
		fmt.Sprintf("Recommend.MA%s", dataInterval),
		fmt.Sprintf("RSI%s", dataInterval),
		fmt.Sprintf("RSI[1]%s", dataInterval),
		fmt.Sprintf("Stoch.K%s", dataInterval),
		fmt.Sprintf("Stoch.D%s", dataInterval),
		fmt.Sprintf("Stoch.K[1]%s", dataInterval),
		fmt.Sprintf("Stoch.D[1]%s", dataInterval),
		fmt.Sprintf("CCI20%s", dataInterval),
		fmt.Sprintf("CCI20[1]%s", dataInterval),
		fmt.Sprintf("ADX%s", dataInterval),
		fmt.Sprintf("ADX+DI%s", dataInterval),
		fmt.Sprintf("ADX-DI%s", dataInterval),
		fmt.Sprintf("ADX+DI[1]%s", dataInterval),
		fmt.Sprintf("ADX-DI[1]%s", dataInterval),
		fmt.Sprintf("AO%s", dataInterval),
		fmt.Sprintf("AO[1]%s", dataInterval),
		fmt.Sprintf("AO[2]%s", dataInterval),
		fmt.Sprintf("MACD.macd%s", dataInterval),
		fmt.Sprintf("MACD.signal%s", dataInterval),
		fmt.Sprintf("Mom%s", dataInterval),
		fmt.Sprintf("Mom[1]%s", dataInterval),
		fmt.Sprintf("Rec.Stoch.RSI%s", dataInterval),
		fmt.Sprintf("Stoch.RSI.K%s", dataInterval),
		fmt.Sprintf("Rec.WR%s", dataInterval),
		fmt.Sprintf("W.R%s", dataInterval),
		fmt.Sprintf("Rec.BBPower%s", dataInterval),
		fmt.Sprintf("BBPower%s", dataInterval),
		fmt.Sprintf("Rec.UO%s", dataInterval),
		fmt.Sprintf("UO%s", dataInterval),
		fmt.Sprintf("EMA10%s", dataInterval),
		fmt.Sprintf("SMA10%s", dataInterval),
		fmt.Sprintf("EMA20%s", dataInterval),
		fmt.Sprintf("SMA20%s", dataInterval),
		fmt.Sprintf("EMA30%s", dataInterval),
		fmt.Sprintf("SMA30%s", dataInterval),
		fmt.Sprintf("EMA50%s", dataInterval),
		fmt.Sprintf("SMA50%s", dataInterval),
		fmt.Sprintf("EMA100%s", dataInterval),
		fmt.Sprintf("SMA100%s", dataInterval),
		fmt.Sprintf("EMA200%s", dataInterval),
		fmt.Sprintf("SMA200%s", dataInterval),
		fmt.Sprintf("Rec.Ichimoku%s", dataInterval),
		fmt.Sprintf("Ichimoku.BLine%s", dataInterval),
		fmt.Sprintf("Rec.VWMA%s", dataInterval),
		fmt.Sprintf("VWMA%s", dataInterval),
		fmt.Sprintf("Rec.HullMA9%s", dataInterval),
		fmt.Sprintf("HullMA9%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.S3%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.S2%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.S1%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.Middle%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.R1%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.R2%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.R3%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.S3%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.S2%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.S1%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.Middle%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.R1%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.R2%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.R3%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.S3%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.S2%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.S1%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.Middle%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.R1%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.R2%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.R3%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.S3%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.S2%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.S1%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.Middle%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.R1%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.R2%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.R3%s", dataInterval),
		fmt.Sprintf("Pivot.M.Demark.S1%s", dataInterval),
		fmt.Sprintf("Pivot.M.Demark.Middle%s", dataInterval),
		fmt.Sprintf("Pivot.M.Demark.R1%s", dataInterval),
		fmt.Sprintf("close%s", dataInterval),
		fmt.Sprintf("high%s", dataInterval),
		fmt.Sprintf("low%s", dataInterval),
	}
}

func (c *Client) newRequest(symbol, dataInterval string) (*http.Request, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("fields", strings.Join(fieldsForInterval(dataInterval), ","))

	reqURL := c.baseURL() + params.Encode()
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) getResponseMap(symbol, dataInterval string) (map[string]float64, error) {
	req, err := c.newRequest(symbol, dataInterval)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	return decodeResponse(res)
}

func decodeResponse(res *http.Response) (map[string]float64, error) {
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status %s: %s", res.Status, strings.TrimSpace(string(jsonData)))
	}

	responseMap := make(map[string]float64)
	if err := json.Unmarshal(jsonData, &responseMap); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return responseMap, nil
}

func (c *Client) baseURL() string {
	var base string
	if c != nil && c.BaseURL != "" {
		base = c.BaseURL
	} else {
		base = defaultScannerURL
	}

	switch {
	case strings.HasSuffix(base, "?"), strings.HasSuffix(base, "&"):
		return base
	case strings.Contains(base, "?"):
		return base + "&"
	default:
		return base + "?"
	}
}

func (c *Client) httpClient() *http.Client {
	if c != nil && c.HTTPClient != nil {
		return c.HTTPClient
	}
	return defaultHTTPClient
}

func (ta *TradingView) populate(responseMap map[string]float64, dataInterval string) {
	ta.populateGlobal(responseMap, dataInterval)
	ta.populateOscillators(responseMap, dataInterval)
	ta.populateMovingAverages(responseMap, dataInterval)
	ta.populatePivots(responseMap, dataInterval)
	ta.populatePrices(responseMap, dataInterval)
}

func (ta *TradingView) populateGlobal(responseMap map[string]float64, dataInterval string) {
	ta.Recommend.Global.Summary = tvComputeRecommend(responseMap[key("Recommend.All%s", dataInterval)])
	ta.Value.Global.Summary = responseMap[key("Recommend.All%s", dataInterval)]

	ta.Recommend.Global.Oscillators = tvComputeRecommend(responseMap[key("Recommend.Other%s", dataInterval)])
	ta.Value.Global.Oscillators = responseMap[key("Recommend.Other%s", dataInterval)]

	ta.Recommend.Global.MA = tvComputeRecommend(responseMap[key("Recommend.MA%s", dataInterval)])
	ta.Value.Global.MA = responseMap[key("Recommend.MA%s", dataInterval)]
}

func (ta *TradingView) populateOscillators(responseMap map[string]float64, dataInterval string) {
	ta.Recommend.Oscillators.RSI = tvRSI(responseMap[key("RSI%s", dataInterval)], responseMap[key("RSI[1]%s", dataInterval)])
	ta.Value.Oscillators.RSI = responseMap[key("RSI%s", dataInterval)]

	ta.Recommend.Oscillators.StochK = tvStoch(responseMap[key("Stoch.K%s", dataInterval)], responseMap[key("Stoch.D%s", dataInterval)], responseMap[key("Stoch.K[1]%s", dataInterval)], responseMap[key("Stoch.D[1]%s", dataInterval)])
	ta.Value.Oscillators.StochK = responseMap[key("Stoch.K%s", dataInterval)]

	ta.Recommend.Oscillators.CCI = tvCCI20(responseMap[key("CCI20%s", dataInterval)], responseMap[key("CCI20[1]%s", dataInterval)])
	ta.Value.Oscillators.CCI = responseMap[key("CCI20%s", dataInterval)]

	ta.Recommend.Oscillators.ADX = tvADX(responseMap[key("ADX%s", dataInterval)], responseMap[key("ADX+DI%s", dataInterval)], responseMap[key("ADX-DI%s", dataInterval)], responseMap[key("ADX+DI[1]%s", dataInterval)], responseMap[key("ADX-DI[1]%s", dataInterval)])
	ta.Value.Oscillators.ADX.Value = responseMap[key("ADX%s", dataInterval)]
	ta.Value.Oscillators.ADX.PlusDI = responseMap[key("ADX+DI%s", dataInterval)]
	ta.Value.Oscillators.ADX.MinusDI = responseMap[key("ADX-DI%s", dataInterval)]
	ta.Value.Oscillators.ADX.PlusDI1 = responseMap[key("ADX+DI[1]%s", dataInterval)]
	ta.Value.Oscillators.ADX.MinusDI1 = responseMap[key("ADX-DI[1]%s", dataInterval)]

	ta.Recommend.Oscillators.AO = tvAO(responseMap[key("AO%s", dataInterval)], responseMap[key("AO[1]%s", dataInterval)], responseMap[key("AO[2]%s", dataInterval)])
	ta.Value.Oscillators.AO.Value = responseMap[key("AO%s", dataInterval)]
	ta.Value.Oscillators.AO.Prev1 = responseMap[key("AO[1]%s", dataInterval)]
	ta.Value.Oscillators.AO.Prev2 = responseMap[key("AO[2]%s", dataInterval)]

	ta.Recommend.Oscillators.Mom = tvMom(responseMap[key("Mom%s", dataInterval)], responseMap[key("Mom[1]%s", dataInterval)])
	ta.Value.Oscillators.Mom = responseMap[key("Mom%s", dataInterval)]

	ta.Recommend.Oscillators.MACD = tvMACD(responseMap[key("MACD.macd%s", dataInterval)], responseMap[key("MACD.signal%s", dataInterval)])
	ta.Value.Oscillators.MACD.Macd = responseMap[key("MACD.macd%s", dataInterval)]
	ta.Value.Oscillators.MACD.Signal = responseMap[key("MACD.signal%s", dataInterval)]

	ta.Recommend.Oscillators.StochRSI = tvSimple(responseMap[key("Rec.Stoch.RSI%s", dataInterval)])
	ta.Value.Oscillators.StochRSI = responseMap[key("Stoch.RSI.K%s", dataInterval)]

	ta.Recommend.Oscillators.WR = tvSimple(responseMap[key("Rec.WR%s", dataInterval)])
	ta.Value.Oscillators.WR = responseMap[key("W.R%s", dataInterval)]

	ta.Recommend.Oscillators.BBP = tvSimple(responseMap[key("Rec.BBPower%s", dataInterval)])
	ta.Value.Oscillators.BBP = responseMap[key("BBPower%s", dataInterval)]

	ta.Recommend.Oscillators.UO = tvSimple(responseMap[key("Rec.UO%s", dataInterval)])
	ta.Value.Oscillators.UO = responseMap[key("UO%s", dataInterval)]
}

func (ta *TradingView) populateMovingAverages(responseMap map[string]float64, dataInterval string) {
	ta.Recommend.MovingAverages.EMA10 = tvMA(responseMap[key("EMA10%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA10 = responseMap[key("EMA10%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA20 = tvMA(responseMap[key("EMA20%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA20 = responseMap[key("EMA20%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA30 = tvMA(responseMap[key("EMA30%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA30 = responseMap[key("EMA30%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA50 = tvMA(responseMap[key("EMA50%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA50 = responseMap[key("EMA50%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA100 = tvMA(responseMap[key("EMA100%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA100 = responseMap[key("EMA100%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA200 = tvMA(responseMap[key("EMA200%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA200 = responseMap[key("EMA200%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA10 = tvMA(responseMap[key("SMA10%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA10 = responseMap[key("SMA10%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA20 = tvMA(responseMap[key("SMA20%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA20 = responseMap[key("SMA20%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA30 = tvMA(responseMap[key("SMA30%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA30 = responseMap[key("SMA30%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA50 = tvMA(responseMap[key("SMA50%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA50 = responseMap[key("SMA50%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA100 = tvMA(responseMap[key("SMA100%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA100 = responseMap[key("SMA100%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA200 = tvMA(responseMap[key("SMA200%s", dataInterval)], responseMap[key("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA200 = responseMap[key("SMA200%s", dataInterval)]

	ta.Recommend.MovingAverages.Ichimoku = tvSimple(responseMap[key("Rec.Ichimoku%s", dataInterval)])
	ta.Value.MovingAverages.Ichimoku = responseMap[key("Ichimoku.BLine%s", dataInterval)]

	ta.Recommend.MovingAverages.VWMA = tvSimple(responseMap[key("Rec.VWMA%s", dataInterval)])
	ta.Value.MovingAverages.VWMA = responseMap[key("VWMA%s", dataInterval)]

	ta.Recommend.MovingAverages.HullMA = tvSimple(responseMap[key("Rec.HullMA9%s", dataInterval)])
	ta.Value.MovingAverages.HullMA = responseMap[key("HullMA9%s", dataInterval)]
}

func (ta *TradingView) populatePivots(responseMap map[string]float64, dataInterval string) {
	ta.Value.Pivots.Classic.Middle = responseMap[key("Pivot.M.Classic.Middle%s", dataInterval)]
	ta.Value.Pivots.Classic.R1 = responseMap[key("Pivot.M.Classic.R1%s", dataInterval)]
	ta.Value.Pivots.Classic.R2 = responseMap[key("Pivot.M.Classic.R2%s", dataInterval)]
	ta.Value.Pivots.Classic.R3 = responseMap[key("Pivot.M.Classic.R3%s", dataInterval)]
	ta.Value.Pivots.Classic.S1 = responseMap[key("Pivot.M.Classic.S1%s", dataInterval)]
	ta.Value.Pivots.Classic.S2 = responseMap[key("Pivot.M.Classic.S2%s", dataInterval)]
	ta.Value.Pivots.Classic.S3 = responseMap[key("Pivot.M.Classic.S3%s", dataInterval)]

	ta.Value.Pivots.Fibonacci.Middle = responseMap[key("Pivot.M.Fibonacci.Middle%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.R1 = responseMap[key("Pivot.M.Fibonacci.R1%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.R2 = responseMap[key("Pivot.M.Fibonacci.R2%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.R3 = responseMap[key("Pivot.M.Fibonacci.R3%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.S1 = responseMap[key("Pivot.M.Fibonacci.S1%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.S2 = responseMap[key("Pivot.M.Fibonacci.S2%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.S3 = responseMap[key("Pivot.M.Fibonacci.S3%s", dataInterval)]

	ta.Value.Pivots.Camarilla.Middle = responseMap[key("Pivot.M.Camarilla.Middle%s", dataInterval)]
	ta.Value.Pivots.Camarilla.R1 = responseMap[key("Pivot.M.Camarilla.R1%s", dataInterval)]
	ta.Value.Pivots.Camarilla.R2 = responseMap[key("Pivot.M.Camarilla.R2%s", dataInterval)]
	ta.Value.Pivots.Camarilla.R3 = responseMap[key("Pivot.M.Camarilla.R3%s", dataInterval)]
	ta.Value.Pivots.Camarilla.S1 = responseMap[key("Pivot.M.Camarilla.S1%s", dataInterval)]
	ta.Value.Pivots.Camarilla.S2 = responseMap[key("Pivot.M.Camarilla.S2%s", dataInterval)]
	ta.Value.Pivots.Camarilla.S3 = responseMap[key("Pivot.M.Camarilla.S3%s", dataInterval)]

	ta.Value.Pivots.Woodie.Middle = responseMap[key("Pivot.M.Woodie.Middle%s", dataInterval)]
	ta.Value.Pivots.Woodie.R1 = responseMap[key("Pivot.M.Woodie.R1%s", dataInterval)]
	ta.Value.Pivots.Woodie.R2 = responseMap[key("Pivot.M.Woodie.R2%s", dataInterval)]
	ta.Value.Pivots.Woodie.R3 = responseMap[key("Pivot.M.Woodie.R3%s", dataInterval)]
	ta.Value.Pivots.Woodie.S1 = responseMap[key("Pivot.M.Woodie.S1%s", dataInterval)]
	ta.Value.Pivots.Woodie.S2 = responseMap[key("Pivot.M.Woodie.S2%s", dataInterval)]
	ta.Value.Pivots.Woodie.S3 = responseMap[key("Pivot.M.Woodie.S3%s", dataInterval)]

	ta.Value.Pivots.Demark.Middle = responseMap[key("Pivot.M.Demark.Middle%s", dataInterval)]
	ta.Value.Pivots.Demark.R1 = responseMap[key("Pivot.M.Demark.R1%s", dataInterval)]
	ta.Value.Pivots.Demark.S1 = responseMap[key("Pivot.M.Demark.S1%s", dataInterval)]
}

func (ta *TradingView) populatePrices(responseMap map[string]float64, dataInterval string) {
	ta.Value.Prices.Close = responseMap[key("close%s", dataInterval)]
	ta.Value.Prices.High = responseMap[key("high%s", dataInterval)]
	ta.Value.Prices.Low = responseMap[key("low%s", dataInterval)]
}

// key formats a response key from an indicator pattern and interval suffix.
func key(indicator, dataInterval string) string {
	return fmt.Sprintf(indicator, dataInterval)
}

// tvComputeRecommend converts TradingView's aggregate score into a public signal.
func tvComputeRecommend(v float64) int {
	switch {
	case v > 0.1 && v <= 0.5:
		return SignalBuy
	case v > 0.5 && v <= 1:
		return SignalStrongBuy
	case v >= -0.1 && v <= 0.1:
		return SignalNeutral
	case v >= -1 && v < -0.5:
		return SignalStrongSell
	case v >= -0.5 && v < -0.1:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// tvRSI converts RSI values into a normalized recommendation signal.
func tvRSI(rsi, rsi1 float64) int {
	switch {
	case rsi < 30 && rsi1 < rsi:
		return SignalBuy
	case rsi > 70 && rsi1 > rsi:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// tvStoch converts stochastic values into a normalized recommendation signal.
func tvStoch(k, d, k1, d1 float64) int {
	switch {
	case k < 20 && d < 20 && k > d && k1 < d1:
		return SignalBuy
	case k > 80 && d > 80 && k < d && k1 > d1:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// tvCCI20 converts CCI values into a normalized recommendation signal.
func tvCCI20(cci20, cci201 float64) int {
	switch {
	case cci20 < -100 && cci20 > cci201:
		return SignalBuy
	case cci20 > 100 && cci20 < cci201:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// tvADX converts ADX values into a normalized recommendation signal.
func tvADX(adx, adxpdi, adxndi, adxpdi1, adxndi1 float64) int {
	switch {
	case adx > 20 && adxpdi1 < adxndi1 && adxpdi > adxndi:
		return SignalBuy
	case adx > 20 && adxpdi1 > adxndi1 && adxpdi < adxndi:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// tvAO converts Awesome Oscillator values into a normalized recommendation signal.
func tvAO(ao, ao1, ao2 float64) int {
	switch {
	case (ao > 0 && ao1 < 0) || (ao > 0 && ao1 > 0 && ao > ao1 && ao2 > ao1):
		return SignalBuy
	case (ao < 0 && ao1 > 0) || (ao < 0 && ao1 < 0 && ao < ao1 && ao2 < ao1):
		return SignalSell
	default:
		return SignalNeutral
	}
}

// tvMom converts momentum values into a normalized recommendation signal.
func tvMom(mom, mom1 float64) int {
	switch {
	case mom > mom1:
		return SignalBuy
	case mom < mom1:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// tvMACD converts MACD values into a normalized recommendation signal.
func tvMACD(macd, s float64) int {
	switch {
	case macd > s:
		return SignalBuy
	case macd < s:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// tvSimple converts a TradingView simple recommendation value into a signal.
func tvSimple(v float64) int {
	switch {
	case v == 1:
		return SignalBuy
	case v == -1:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// tvMA converts a moving average and close price into a signal.
func tvMA(ma, close float64) int {
	switch {
	case ma < close:
		return SignalBuy
	case ma > close:
		return SignalSell
	default:
		return SignalNeutral
	}
}
