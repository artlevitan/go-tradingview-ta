// Copyright 2022-2023. All rights reserved.
// https://github.com/artlevitan/go-tradingview-ta
// v1.2.2

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
	Interval1min   string = "1min"
	Interval5min   string = "5min"
	Interval15min  string = "15min"
	Interval30min  string = "30min"
	Interval1hour  string = "1hour"
	Interval2hour  string = "2hour"
	Interval4hour  string = "4hour"
	Interval1day   string = "1day"
	Interval1week  string = "1week"
	Interval1month string = "1month"

	SignalStrongBuy  int = 2  // STRONG_BUY
	SignalBuy        int = 1  // BUY
	SignalNeutral    int = 0  // NEUTRAL
	SignalSell       int = -1 // SELL
	SignalStrongSell int = -2 // STRONG_SELL
)

// TradingView Payload Data
type TradingView struct {
	Recommend struct {
		Summary     int // Summary
		Oscillators int // Oscillators
		MA          int // Moving Averages
	}
	Oscillators struct {
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
	MovingAverages struct {
		EMA10    int // Exponential Moving Average (10)
		SMA10    int // Simple Moving Average (10)
		EMA20    int // Exponential Moving Average (20)
		SMA20    int // Simple Moving Average (20)
		EMA30    int // Exponential Moving Average (30)
		SMA30    int // Simple Moving Average (30)
		EMA50    int // Exponential Moving Average (50)
		SMA50    int // Simple Moving Average (50)
		EMA100   int // Exponential Moving Average (100)
		SMA100   int // Simple Moving Average (100)
		EMA200   int // Exponential Moving Average (200)
		SMA200   int // Simple Moving Average (200)
		Ichimoku int // Ichimoku Base Line (9, 26, 52, 26)
		VWMA     int // Volume Weighted Moving Average (20)
		HullMA   int // Hull Moving Average (9)
	}
}

// Get - format TradingView's Scanner Post Data
//
// symbol string – Name of EXCHANGE:SYMBOL (ex: "BINANCE:BTCUSDT" or "BINANCE:ETHUSDT")
//
// interval string - Interval / Timeframe
func (t *TradingView) Get(symbol string, interval string) error {
	// Parameters validation
	if strings.Count(symbol, ":") != 1 {
		return errors.New("symbol parameter is not valid")
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
		return errors.New(err.Error())
	}

	// Headers
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return errors.New(err.Error())
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New(err.Error())
	}

	// Indicators data
	type Indicators struct {
		TotalCount int `json:"totalCount"`
		Data       []struct {
			S string    `json:"s"`
			D []float64 `json:"d"`
		} `json:"data"`
	}

	inds := Indicators{}
	err = json.Unmarshal(jsonData, &inds)
	if err != nil {
		return errors.New(err.Error())
	}

	// Data not received
	if inds.TotalCount == 0 {
		return errors.New("data not received")
	}

	// Recommendations
	t.Recommend.Summary = tvComputeRecommend(inds.Data[0].D[0])
	t.Recommend.Oscillators = tvComputeRecommend(inds.Data[0].D[1])
	t.Recommend.MA = tvComputeRecommend(inds.Data[0].D[2])

	// Oscillators
	// Relative Strength Index (14)
	t.Oscillators.RSI = tvRsi(inds.Data[0].D[3], inds.Data[0].D[4])

	// Stochastic %K (14, 3, 3)
	t.Oscillators.StochK = tvStoch(inds.Data[0].D[5], inds.Data[0].D[6], inds.Data[0].D[7], inds.Data[0].D[7])

	// Commodity Channel Index (20)
	t.Oscillators.CCI = tvCci20(inds.Data[0].D[9], inds.Data[0].D[10])

	// Average Directional Index (14)
	t.Oscillators.ADX = tvAdx(inds.Data[0].D[11], inds.Data[0].D[12], inds.Data[0].D[13], inds.Data[0].D[14], inds.Data[0].D[15])

	// Awesome Oscillator
	t.Oscillators.AO = tvAo(inds.Data[0].D[16], inds.Data[0].D[17], inds.Data[0].D[86])

	// Momentum (10)
	t.Oscillators.Mom = tvMom(inds.Data[0].D[18], inds.Data[0].D[19])

	// MACD Level (12, 26)
	t.Oscillators.MACD = tvMacd(inds.Data[0].D[20], inds.Data[0].D[21])

	// Stochastic RSI Fast (3, 3, 14, 14)
	t.Oscillators.StochRSI = tvSimple(inds.Data[0].D[22])

	// Williams Percent Range (14)
	t.Oscillators.WR = tvSimple(inds.Data[0].D[24])

	// Bull Bear Power
	t.Oscillators.BBP = tvSimple(inds.Data[0].D[26])

	// Ultimate Oscillator (7, 14, 28)
	t.Oscillators.UO = tvSimple(inds.Data[0].D[28])

	// Moving Averages
	for i := 33; i < 45; i++ {
		switch i {
		case 33:
			// Exponential Moving Average (10)
			t.MovingAverages.EMA10 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 34:
			// Simple Moving Average (10)
			t.MovingAverages.SMA10 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 35:
			// Exponential Moving Average (20)
			t.MovingAverages.EMA20 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 36:
			// Simple Moving Average (20)
			t.MovingAverages.SMA20 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 37:
			// Exponential Moving Average (30)
			t.MovingAverages.EMA30 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 38:
			// Simple Moving Average (30)
			t.MovingAverages.SMA30 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 39:
			// Exponential Moving Average (50)
			t.MovingAverages.EMA50 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 40:
			// Simple Moving Average (50)
			t.MovingAverages.SMA50 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 41:
			// Exponential Moving Average (100)
			t.MovingAverages.EMA100 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 42:
			// Simple Moving Average (100)
			t.MovingAverages.SMA100 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 43:
			// Exponential Moving Average (200)
			t.MovingAverages.EMA200 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		case 44:
			// Simple Moving Average (200)
			t.MovingAverages.SMA200 = tvMa(inds.Data[0].D[i], inds.Data[0].D[30])
		}
	}

	// Ichimoku Base Line (9, 26, 52, 26)
	t.MovingAverages.Ichimoku = tvSimple(inds.Data[0].D[45])

	// Volume Weighted Moving Average (20)
	t.MovingAverages.VWMA = tvSimple(inds.Data[0].D[47])

	// Hull Moving Average (9)
	t.MovingAverages.HullMA = tvSimple(inds.Data[0].D[49])

	return nil
}

// tvComputeRecommend - Compute Recommend
func tvComputeRecommend(v float64) int {
	switch {
	case v > 0.1 && v <= 0.5:
		return SignalBuy // BUY
	case v > 0.5 && v <= 1:
		return SignalStrongBuy // STRONG_BUY
	case v >= -0.1 && v <= 0.1:
		return SignalNeutral // NEUTRAL
	case v >= -1 && v < -0.5:
		return SignalStrongSell // STRONG_SELL
	case v >= -0.5 && v < -0.1:
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}

// tvRsi - Compute Relative Strength Index
func tvRsi(rsi, rsi1 float64) int {
	switch {
	case rsi < 30 && rsi1 < rsi:
		return SignalBuy // BUY
	case rsi > 70 && rsi1 > rsi:
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}

// tvStoch - Compute Stochastic
func tvStoch(k, d, k1, d1 float64) int {
	switch {
	case k < 20 && d < 20 && k > d && k1 < d1:
		return SignalBuy // BUY
	case k > 80 && d > 80 && k < d && k1 > d1:
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}

// tvCci20 - Compute Commodity Channel Index 20
func tvCci20(cci20, cci201 float64) int {
	switch {
	case cci20 < -100 && cci20 > cci201:
		return SignalBuy // BUY
	case cci20 > 100 && cci20 < cci201:
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}

// tvAdx - Compute Average Directional Index
func tvAdx(adx, adxpdi, adxndi, adxpdi1, adxndi1 float64) int {
	switch {
	case adx > 20 && adxpdi1 < adxndi1 && adxpdi > adxndi:
		return SignalBuy // BUY
	case adx > 20 && adxpdi1 > adxndi1 && adxpdi < adxndi:
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}

// tvAo - Compute Awesome Oscillator
func tvAo(ao, ao1, ao2 float64) int {
	switch {
	case (ao > 0 && ao1 < 0) || (ao > 0 && ao1 > 0 && ao > ao1 && ao2 > ao1):
		return SignalBuy // BUY
	case (ao < 0 && ao1 > 0) || (ao < 0 && ao1 < 0 && ao < ao1 && ao2 < ao1):
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}

// tvMom - Compute Momentum
func tvMom(mom, mom1 float64) int {
	switch {
	case mom > mom1:
		return SignalBuy // BUY
	case mom < mom1:
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}

// tvMacd - Compute Moving Average Convergence/Divergence
func tvMacd(macd, s float64) int {
	switch {
	case macd > s:
		return SignalBuy // BUY
	case macd < s:
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}

// tvSimple - Compute Simple
func tvSimple(v float64) int {
	switch {
	case v == 1:
		return SignalBuy // BUY
	case v == -1:
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}

// tvMa - Compute Moving Average
func tvMa(ma, close float64) int {
	switch {
	case ma < close:
		return SignalBuy // BUY
	case ma > close:
		return SignalSell // SELL
	default:
		return SignalNeutral // NEUTRAL
	}
}
