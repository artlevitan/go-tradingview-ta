// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import "time"

// Supported interval values for Get and GetContext.
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
	// SignalUnavailable means that one or more required values are null.
	// It is not a recommendation and must not be used as a numeric rating.
	SignalUnavailable = 3
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
// Its zero value is ready to use. A TradingView must not be read or written
// concurrently with a call that updates it.
type TradingView struct {
	// Recommend contains normalized recommendation signals.
	Recommend Recommendations
	// Value contains raw numeric indicator and price values.
	// Unavailable values are zero; use HasValue or Raw to distinguish them.
	Value Values
	// Raw preserves numeric scanner values, including nulls. Keys are scanner
	// field names without interval suffixes, such as "EMA200" or "RSI[1]".
	Raw *RawValues
	// Metadata describes the instrument, candle, and data update mode.
	Metadata Metadata
}

// Metadata describes a successful scanner snapshot.
type Metadata struct {
	Symbol     string    // Exchange-qualified symbol, such as "BINANCE:BTCUSDT".
	Interval   string    // Normalized interval; an empty request is recorded as Interval1Day.
	FetchedAt  time.Time // Time when the response was successfully decoded (UTC).
	BarTime    time.Time // Opening time of the requested candle (UTC), or zero if unavailable.
	UpdateMode string    // Scanner mode, e.g. "streaming" or "delayed_streaming_900".
}

// HasValue reports whether a scanner field has a numeric value in this snapshot.
// Field names do not include interval suffixes, for example "EMA200".
func (ta *TradingView) HasValue(field string) bool {
	if ta == nil {
		return false
	}
	_, ok := ta.Raw.Lookup(field)
	return ok
}

// Recommendations groups normalized signals, or SignalUnavailable when the
// inputs needed for an indicator's recommendation are not available.
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

// Values groups numeric values returned by TradingView.
// Unavailable fields contain zero; use TradingView.HasValue or RawValues.Lookup
// to distinguish them from genuine zero values.
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
	Classic   ClassicPivotLevels   // Classic support, resistance, and central pivot levels.
	Fibonacci FibonacciPivotLevels // Fibonacci support, resistance, and central pivot levels.
	Camarilla CamarillaPivotLevels // Camarilla support, resistance, and central pivot levels.
	Woodie    WoodiePivotLevels    // Woodie support, resistance, and central pivot levels.
	Demark    DemarkPivotLevels    // DeMark support, resistance, and central pivot levels.
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

// PriceValues stores the requested candle's prices returned by TradingView.
// The candle may still be open, in which case Close is its latest price.
type PriceValues struct {
	Close float64 // Closing price
	High  float64 // Highest price
	Low   float64 // Lowest price
}
