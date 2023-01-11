// Copyright 2022-2023. All rights reserved.
// https://github.com/artlevitan/go-tradingview-ta
// 1.0.0
// License: MIT

package tradingview

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// Intervals
	Interval1min   = "1min"
	Interval5min   = "5min"
	Interval15min  = "15min"
	Interval30min  = "30min"
	Interval1hour  = "1hour"
	Interval2hour  = "2hour"
	Interval4hour  = "4hour"
	Interval1day   = "1day"
	Interval1week  = "1week"
	Interval1month = "1month"

	// Result
	SignalStrongSell = -2 // STRONG_SELL
	SignalSell       = -1 // SELL
	SignalNeutral    = 0  // NEUTRAL
	SignalBuy        = 1  // BUY
	SignalStrongBuy  = 2  // STRONG_BUY
)

// TradingViewData - format TradingView's Scanner Post Data.
//
// symbols – Name of EXCHANGE:SYMBOL (ex: "BINANCE:BTCUSDT" or "BINANCE:ETHUSDT")
//
// interval - Interval / Timeframe
func TradingViewData(symbol string, interval string) (map[string]int, error) {
	result := map[string]int{}

	// Parameters validation
	if strings.Count(symbol, ":") != 1 {
		return result, errors.New("symbol parameter is not valid")
	}

	var dataInterval string
	switch interval {
	case Interval1min:
		//  1 Minute
		dataInterval = "1"
	case Interval5min:
		//  5 Minutes
		dataInterval = "5"
	case Interval15min:
		//  15 Minutes
		dataInterval = "15"
	case Interval30min:
		//  30 Minutes
		dataInterval = "30"
	case Interval1hour:
		//  1 Hour
		dataInterval = "60"
	case Interval2hour:
		//  2 Hours
		dataInterval = "120"
	case Interval4hour:
		//  4 Hour
		dataInterval = "240"
	case Interval1week:
		//  1 Week
		dataInterval = "1W"
	case Interval1month:
		//  1 Month
		dataInterval = "1M"
	default: // Default 1 day
		dataInterval = Interval1day
	}

	// Request preparation
	type Request struct {
		Symbols struct {
			Tickers []string `json:"tickers"`
		} `json:"symbols"`
		Columns []string `json:"columns"`
	}
	data := Request{}
	data.Symbols.Tickers = []string{symbol}
	data.Columns = []string{
		fmt.Sprintf("Recommend.All|%s", dataInterval),
		fmt.Sprintf("Recommend.Other|%s", dataInterval),
		fmt.Sprintf("Recommend.MA|%s", dataInterval),
		fmt.Sprintf("RSI|%s", dataInterval),
		fmt.Sprintf("RSI[1]|%s", dataInterval),
		fmt.Sprintf("Stoch.K|%s", dataInterval),
		fmt.Sprintf("Stoch.D|%s", dataInterval),
		fmt.Sprintf("Stoch.K[1]|%s", dataInterval),
		fmt.Sprintf("Stoch.D[1]|%s", dataInterval),
		fmt.Sprintf("CCI20|%s", dataInterval),
		fmt.Sprintf("CCI20[1]|%s", dataInterval),
		fmt.Sprintf("ADX|%s", dataInterval),
		fmt.Sprintf("ADX+DI|%s", dataInterval),
		fmt.Sprintf("ADX-DI|%s", dataInterval),
		fmt.Sprintf("ADX+DI[1]|%s", dataInterval),
		fmt.Sprintf("ADX-DI[1]|%s", dataInterval),
		fmt.Sprintf("AO|%s", dataInterval),
		fmt.Sprintf("AO[1]|%s", dataInterval),
		fmt.Sprintf("Mom|%s", dataInterval),
		fmt.Sprintf("Mom[1]|%s", dataInterval),
		fmt.Sprintf("MACD.macd|%s", dataInterval),
		fmt.Sprintf("MACD.signal|%s", dataInterval),
		fmt.Sprintf("Rec.Stoch.RSI|%s", dataInterval),
		fmt.Sprintf("Stoch.RSI.K|%s", dataInterval),
		fmt.Sprintf("Rec.WR|%s", dataInterval),
		fmt.Sprintf("W.R|%s", dataInterval),
		fmt.Sprintf("Rec.BBPower|%s", dataInterval),
		fmt.Sprintf("BBPower|%s", dataInterval),
		fmt.Sprintf("Rec.UO|%s", dataInterval),
		fmt.Sprintf("UO|%s", dataInterval),
		fmt.Sprintf("close|%s", dataInterval),
		fmt.Sprintf("EMA5|%s", dataInterval),
		fmt.Sprintf("SMA5|%s", dataInterval),
		fmt.Sprintf("EMA10|%s", dataInterval),
		fmt.Sprintf("SMA10|%s", dataInterval),
		fmt.Sprintf("EMA20|%s", dataInterval),
		fmt.Sprintf("SMA20|%s", dataInterval),
		fmt.Sprintf("EMA30|%s", dataInterval),
		fmt.Sprintf("SMA30|%s", dataInterval),
		fmt.Sprintf("EMA50|%s", dataInterval),
		fmt.Sprintf("SMA50|%s", dataInterval),
		fmt.Sprintf("EMA100|%s", dataInterval),
		fmt.Sprintf("SMA100|%s", dataInterval),
		fmt.Sprintf("EMA200|%s", dataInterval),
		fmt.Sprintf("SMA200|%s", dataInterval),
		fmt.Sprintf("Rec.Ichimoku|%s", dataInterval),
		fmt.Sprintf("Ichimoku.BLine|%s", dataInterval),
		fmt.Sprintf("Rec.VWMA|%s", dataInterval),
		fmt.Sprintf("VWMA|%s", dataInterval),
		fmt.Sprintf("Rec.HullMA9|%s", dataInterval),
		fmt.Sprintf("HullMA9|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.S3|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.S2|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.S1|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.Middle|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.R1|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.R2|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Classic.R3|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.S3|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.S2|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.S1|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.Middle|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.R1|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.R2|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Fibonacci.R3|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.S3|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.S2|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.S1|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.Middle|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.R1|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.R2|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Camarilla.R3|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.S3|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.S2|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.S1|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.Middle|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.R1|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.R2|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Woodie.R3|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Demark.S1|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Demark.Middle|%s", dataInterval),
		fmt.Sprintf("Pivot.M.Demark.R1|%s", dataInterval),
		fmt.Sprintf("open|%s", dataInterval),
		fmt.Sprintf("P.SAR|%s", dataInterval),
		fmt.Sprintf("BB.lower|%s", dataInterval),
		fmt.Sprintf("BB.upper|%s", dataInterval),
		fmt.Sprintf("AO[2]|%s", dataInterval),
		fmt.Sprintf("volume|%s", dataInterval),
		fmt.Sprintf("change|%s", dataInterval),
		fmt.Sprintf("low|%s", dataInterval),
		fmt.Sprintf("high|%s", dataInterval),
	}

	bytes, _ := json.Marshal(data)
	payload := string(bytes)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://scanner.tradingview.com/crypto/scan", strings.NewReader(payload))
	if err != nil {
		return result, errors.New(err.Error())
	}

	// Headers
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return result, errors.New(err.Error())
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return result, errors.New(err.Error())
	}

	type TradingView struct {
		TotalCount int `json:"totalCount"`
		Data       []struct {
			S string    `json:"s"`
			D []float64 `json:"d"`
		} `json:"data"`
	}

	indicators := TradingView{}
	err = json.Unmarshal(jsonData, &indicators)
	if err != nil {
		return result, errors.New(err.Error())
	}

	// Data not received
	if indicators.TotalCount == 0 {
		return result, errors.New("data not received")
	}

	// RECOMMENDATIONS
	result["recommend_summary"] = tvComputerecommend(indicators.Data[0].D[0])
	result["recommend_oscillators"] = tvComputerecommend(indicators.Data[0].D[1])
	result["recommend_moving_averages"] = tvComputerecommend(indicators.Data[0].D[2])

	// OSCILLATORS
	// RSI (14)
	result["computed_oscillators_RSI"] = tvRsi(indicators.Data[0].D[3], indicators.Data[0].D[4])

	// Stoch %K
	result["computed_oscillators_STOCHK"] = tvStoch(indicators.Data[0].D[5], indicators.Data[0].D[6], indicators.Data[0].D[7], indicators.Data[0].D[7])

	// CCI (20)
	result["computed_oscillators_CCI"] = tvCci20(indicators.Data[0].D[9], indicators.Data[0].D[10])

	// ADX (14)
	result["computed_oscillators_ADX"] = tvAdx(indicators.Data[0].D[11], indicators.Data[0].D[12], indicators.Data[0].D[13], indicators.Data[0].D[14], indicators.Data[0].D[15])

	// AO
	result["computed_oscillators_AO"] = tvAo(indicators.Data[0].D[16], indicators.Data[0].D[17], indicators.Data[0].D[86])

	// Mom (10)
	result["computed_oscillators_Mom"] = tvMom(indicators.Data[0].D[18], indicators.Data[0].D[19])

	// MACD
	result["computed_oscillators_MACD"] = tvMacd(indicators.Data[0].D[20], indicators.Data[0].D[21])

	// Stoch RSI
	result["computed_oscillators_STOCHRSI"] = tvSimple(indicators.Data[0].D[22])

	// W%R
	result["computed_oscillators_WR"] = tvSimple(indicators.Data[0].D[24])

	// BBP
	result["computed_oscillators_BBP"] = tvSimple(indicators.Data[0].D[26])

	// UO
	result["computed_oscillators_UO"] = tvSimple(indicators.Data[0].D[28])

	// MOVING AVERAGES
	maList := []string{"EMA10", "SMA10", "EMA20", "SMA20", "EMA30", "SMA30", "EMA50", "SMA50", "EMA100", "SMA100", "EMA200", "SMA200"}
	maListCounter := 0
	for i := 33; i < 45; i++ {
		result["computed_ma_"+maList[maListCounter]] = tvMa(indicators.Data[0].D[i], indicators.Data[0].D[30])
		maListCounter++
	}

	// ICHIMOKU
	result["computed_ma_Ichimoku"] = tvSimple(indicators.Data[0].D[45])

	// VWMA
	result["computed_ma_VWMA"] = tvSimple(indicators.Data[0].D[47])

	// HullMA
	result["computed_ma_HullMA"] = tvSimple(indicators.Data[0].D[49])

	return result, nil
}

// Compute Recommend
func tvComputerecommend(v float64) int {
	switch {
	case v >= -1 && v < -0.5:
		return SignalStrongSell // strong_sell
	case v >= -0.5 && v < -0.1:
		return SignalSell
	case v >= -0.1 && v <= 0.1:
		return SignalNeutral // SignalNeutral
	case v > 0.1 && v <= 0.5:
		return SignalBuy
	case v > 0.5 && v <= 1:
		return SignalStrongBuy // strong_buy
	default:
		return SignalNeutral
	}
}

// Compute Relative Strength Index
func tvRsi(rsi, rsi1 float64) int {
	switch {
	case rsi < 30 && rsi1 < rsi:
		return SignalBuy
	case rsi > 70 && rsi1 > rsi:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// Compute Stochastic
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

// Compute Commodity Channel Index 20
func tvCci20(cci20, cci201 float64) int {
	switch {
	case cci20 < -100 && cci20 > cci201:
		return SignalBuy
	case cci20 > 100 && cci20 < cci201:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// Compute Average Directional Index
func tvAdx(adx, adxpdi, adxndi, adxpdi1, adxndi1 float64) int {
	switch {
	case adx > 20 && adxpdi1 < adxndi1 && adxpdi > adxndi:
		return SignalBuy
	case adx > 20 && adxpdi1 > adxndi1 && adxpdi < adxndi:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// Compute Awesome Oscillator
func tvAo(ao, ao1, ao2 float64) int {
	switch {
	case (ao > 0 && ao1 < 0) || (ao > 0 && ao1 > 0 && ao > ao1 && ao2 > ao1):
		return SignalBuy
	case (ao < 0 && ao1 > 0) || (ao < 0 && ao1 < 0 && ao < ao1 && ao2 < ao1):
		return SignalSell
	default:
		return SignalNeutral
	}
}

// Compute Momentum
func tvMom(mom, mom1 float64) int {
	switch {
	case mom < mom1:
		return SignalSell
	case mom > mom1:
		return SignalBuy
	default:
		return SignalNeutral
	}
}

// Compute Moving Average Convergence/Divergence
func tvMacd(macd, s float64) int {
	switch {
	case macd > s:
		return SignalBuy
	case macd < s:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// Compute Simple
func tvSimple(v float64) int {
	switch {
	case v == -1:
		return SignalSell
	case v == 1:
		return SignalBuy
	default:
		return SignalNeutral
	}
}

// Compute Moving Average
func tvMa(ma, close float64) int {
	switch {
	case ma < close:
		return SignalBuy
	case ma > close:
		return SignalSell
	default:
		return SignalNeutral
	}
}
