// Copyright 2022-2024. All rights reserved.
// https://github.com/artlevitan/go-tradingview-ta
// v1.3.1

package tradingview

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	Interval1Min   string = "1"   // 1 minute
	Interval5Min   string = "5"   // 5 minutes
	Interval15Min  string = "15"  // 15 minutes
	Interval30Min  string = "30"  // 30 minutes
	Interval1Hour  string = "60"  // 1 hour
	Interval2Hour  string = "120" // 2 hours
	Interval4Hour  string = "240" // 4 hours
	Interval1Day   string = "1D"  // 1 day
	Interval1Week  string = "1W"  // 1 week
	Interval1Month string = "1M"  // 1 month

	SignalStrongBuy  int = 2  // STRONG_BUY
	SignalBuy        int = 1  // BUY
	SignalNeutral    int = 0  // NEUTRAL
	SignalSell       int = -1 // SELL
	SignalStrongSell int = -2 // STRONG_SELL

	// Deprecated
	Interval1min   string = Interval1Min
	Interval5min   string = Interval5Min
	Interval15min  string = Interval15Min
	Interval30min  string = Interval30Min
	Interval1hour  string = Interval1Hour
	Interval2hour  string = Interval2Hour
	Interval4hour  string = Interval4Hour
	Interval1day   string = Interval1Day
	Interval1week  string = Interval1Week
	Interval1month string = Interval1Month
)

// TradingView Payload Data
type TradingView struct {
	Recommend struct {
		Global struct {
			Summary     int // Summary recommendation
			Oscillators int // Oscillators recommendation
			MA          int // Moving Averages recommendation
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
	}
	Value struct {
		Global struct {
			Summary     float64 // Summary recommendation
			Oscillators float64 // Oscillators recommendation
			MA          float64 // Moving Averages recommendation
		}
		Oscillators struct {
			RSI    float64  // Relative Strength Index (14)
			StochK float64  // Stochastic %K (14, 3, 3)
			CCI    float64  // Commodity Channel Index (20)
			ADX    struct { // Average Directional Index (14)
				Value    float64 // ADX Value
				PlusDI   float64 // ADX+DI
				MinusDI  float64 // ADX-DI
				PlusDI1  float64 // ADX+DI[1]
				MinusDI1 float64 // ADX-DI[1]
			}
			AO struct { // Awesome Oscillator
				Value float64 // AO current value
				Prev1 float64 // AO[1]
				Prev2 float64 // AO[2]
			}
			Mom  float64  // Momentum (10)
			MACD struct { // MACD Level (12, 26)
				Macd   float64 // MACD line
				Signal float64 // Signal line
			}
			StochRSI float64 // Stochastic RSI Fast (3, 3, 14, 14)
			WR       float64 // Williams Percent Range (14)
			BBP      float64 // Bull Bear Power
			UO       float64 // Ultimate Oscillator (7, 14, 28)
		}
		MovingAverages struct {
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
		Pivots struct {
			Classic struct {
				Middle float64 // Classic Pivot Middle (Pivot.M.Classic.Middle)
				R1     float64 // Resistance 1 (Pivot.M.Classic.R1)
				R2     float64 // Resistance 2 (Pivot.M.Classic.R2)
				R3     float64 // Resistance 3 (Pivot.M.Classic.R3)
				S1     float64 // Support 1 (Pivot.M.Classic.S1)
				S2     float64 // Support 2 (Pivot.M.Classic.S2)
				S3     float64 // Support 3 (Pivot.M.Classic.S3)
			}
			Fibonacci struct {
				Middle float64 // Fibonacci Pivot Middle (Pivot.M.Fibonacci.Middle)
				R1     float64 // Resistance 1 (Pivot.M.Fibonacci.R1)
				R2     float64 // Resistance 2 (Pivot.M.Fibonacci.R2)
				R3     float64 // Resistance 3 (Pivot.M.Fibonacci.R3)
				S1     float64 // Support 1 (Pivot.M.Fibonacci.S1)
				S2     float64 // Support 2 (Pivot.M.Fibonacci.S2)
				S3     float64 // Support 3 (Pivot.M.Fibonacci.S3)
			}
			Camarilla struct {
				Middle float64 // Camarilla Pivot Middle (Pivot.M.Camarilla.Middle)
				R1     float64 // Resistance 1 (Pivot.M.Camarilla.R1)
				R2     float64 // Resistance 2 (Pivot.M.Camarilla.R2)
				R3     float64 // Resistance 3 (Pivot.M.Camarilla.R3)
				S1     float64 // Support 1 (Pivot.M.Camarilla.S1)
				S2     float64 // Support 2 (Pivot.M.Camarilla.S2)
				S3     float64 // Support 3 (Pivot.M.Camarilla.S3)
			}
			Woodie struct {
				Middle float64 // Woodie Pivot Middle (Pivot.M.Woodie.Middle)
				R1     float64 // Resistance 1 (Pivot.M.Woodie.R1)
				R2     float64 // Resistance 2 (Pivot.M.Woodie.R2)
				R3     float64 // Resistance 3 (Pivot.M.Woodie.R3)
				S1     float64 // Support 1 (Pivot.M.Woodie.S1)
				S2     float64 // Support 2 (Pivot.M.Woodie.S2)
				S3     float64 // Support 3 (Pivot.M.Woodie.S3)
			}
			Demark struct {
				Middle float64 // Demark Pivot Middle (Pivot.M.Demark.Middle)
				R1     float64 // Resistance 1 (Pivot.M.Demark.R1)
				S1     float64 // Support 1 (Pivot.M.Demark.S1)
			}
		}
		Prices struct {
			Close float64 // Closing price
			High  float64 // Highest price
			Low   float64 // Lowest price
		}
	}
}

// Get formats and sends a GET request to TradingView's scanner to retrieve market data
// for a given symbol and timeframe.
//
// The symbol should be in the format "EXCHANGE:SYMBOL" (e.g., "BINANCE:BTCUSDT").
// The interval parameter defines the timeframe (e.g., "1min", "5min", "1hour", etc.).
//
// Parameters:
//
//	symbol  - string: the exchange and symbol in the format "EXCHANGE:SYMBOL".
//	interval - string: the timeframe/interval for the data.
//
// Returns:
//
//	error - returns an error if the request fails or the symbol/interval is invalid.
func (ta *TradingView) Get(symbol string, interval string) error {
	// Validate symbol parameter to ensure it has the correct format EXCHANGE:SYMBOL
	if strings.Count(symbol, ":") != 1 {
		return errors.New("symbol parameter is not valid")
	}

	// Map interval input to appropriate TradingView format
	var dataInterval string
	switch interval {
	case Interval1min: // 1 minute
		dataInterval = "|1"
	case Interval5min: // 5 minutes
		dataInterval = "|5"
	case Interval15min: // 15 minutes
		dataInterval = "|15"
	case Interval30min: // 30 minutes
		dataInterval = "|30"
	case Interval1hour: // 1 hour
		dataInterval = "|60"
	case Interval2hour: // 2 hours
		dataInterval = "|120"
	case Interval4hour: // 4 hours
		dataInterval = "|240"
	case Interval1day: // 1 day
		dataInterval = ""
	case Interval1week: // 1 week
		dataInterval = "|1W"
	case Interval1month: // 1 month
		dataInterval = "|1M"
	default: // 1 day
		dataInterval = ""
	}

	// Construct the fields array which includes technical indicators for the specified interval
	fields := []string{
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

	// Build the URL for GET request by encoding the parameters
	baseURL := "https://scanner.tradingview.com/symbol?"
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("fields", strings.Join(fields, ","))

	// Full URL for GET request
	reqURL := baseURL + params.Encode()

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return errors.New("failed to create request: " + err.Error())
	}

	// Set request headers
	req.Header.Add("Content-Type", "application/json")

	// Execute the GET request
	res, err := client.Do(req)
	if err != nil {
		return errors.New("failed to send request: " + err.Error())
	}
	defer res.Body.Close()

	// Read the response body
	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("failed to read response: " + err.Error())
	}

	// Unmarshal the JSON data into a map
	responseMap := make(map[string]float64)
	err = json.Unmarshal(jsonData, &responseMap)
	if err != nil {
		return errors.New("failed to parse response: " + err.Error())
	}

	// Recommendations
	// Summary recommendation
	ta.Recommend.Global.Summary = tvComputeRecommend(responseMap[fmt.Sprintf("Recommend.All%s", dataInterval)])
	ta.Value.Global.Summary = responseMap[fmt.Sprintf("Recommend.All%s", dataInterval)]

	// Oscillators recommendation
	ta.Recommend.Global.Oscillators = tvComputeRecommend(responseMap[fmt.Sprintf("Recommend.Other%s", dataInterval)])
	ta.Value.Global.Oscillators = responseMap[fmt.Sprintf("Recommend.Other%s", dataInterval)]

	// Moving Averages recommendation
	ta.Recommend.Global.MA = tvComputeRecommend(responseMap[fmt.Sprintf("Recommend.MA%s", dataInterval)])
	ta.Value.Global.MA = responseMap[fmt.Sprintf("Recommend.MA%s", dataInterval)]

	// Oscillators
	// Relative Strength Index (14)
	ta.Recommend.Oscillators.RSI = tvRsi(responseMap[fmt.Sprintf("RSI%s", dataInterval)], responseMap[fmt.Sprintf("RSI[1]%s", dataInterval)])
	ta.Value.Oscillators.RSI = responseMap[fmt.Sprintf("RSI%s", dataInterval)]

	// Stochastic %K (14, 3, 3)
	ta.Recommend.Oscillators.StochK = tvStoch(responseMap[fmt.Sprintf("Stoch.K%s", dataInterval)], responseMap[fmt.Sprintf("Stoch.D%s", dataInterval)], responseMap[fmt.Sprintf("Stoch.K[1]%s", dataInterval)], responseMap[fmt.Sprintf("Stoch.D[1]%s", dataInterval)])
	ta.Value.Oscillators.StochK = responseMap[fmt.Sprintf("Stoch.K%s", dataInterval)]

	// Commodity Channel Index (20)
	ta.Recommend.Oscillators.CCI = tvCci20(responseMap[fmt.Sprintf("CCI20%s", dataInterval)], responseMap[fmt.Sprintf("CCI20[1]%s", dataInterval)])
	ta.Value.Oscillators.CCI = responseMap[fmt.Sprintf("CCI20%s", dataInterval)]

	// Average Directional Index (14)
	ta.Recommend.Oscillators.ADX = tvAdx(responseMap[fmt.Sprintf("ADX%s", dataInterval)], responseMap[fmt.Sprintf("ADX+DI%s", dataInterval)], responseMap[fmt.Sprintf("ADX+DI%s", dataInterval)], responseMap[fmt.Sprintf("ADX+DI[1]%s", dataInterval)], responseMap[fmt.Sprintf("ADX-DI[1]%s", dataInterval)])
	ta.Value.Oscillators.ADX.Value = responseMap[fmt.Sprintf("ADX%s", dataInterval)]          // ADX Value
	ta.Value.Oscillators.ADX.PlusDI = responseMap[fmt.Sprintf("ADX+DI%s", dataInterval)]      // ADX +DI
	ta.Value.Oscillators.ADX.MinusDI = responseMap[fmt.Sprintf("ADX-DI%s", dataInterval)]     // ADX -DI
	ta.Value.Oscillators.ADX.PlusDI1 = responseMap[fmt.Sprintf("ADX+DI[1]%s", dataInterval)]  // ADX +DI[1]
	ta.Value.Oscillators.ADX.MinusDI1 = responseMap[fmt.Sprintf("ADX-DI[1]%s", dataInterval)] // ADX -DI[1]

	// Awesome Oscillator
	ta.Recommend.Oscillators.AO = tvAo(responseMap[fmt.Sprintf("AO%s", dataInterval)], responseMap[fmt.Sprintf("AO[1]%s", dataInterval)], responseMap[fmt.Sprintf("AO[2]%s", dataInterval)])
	ta.Value.Oscillators.AO.Value = responseMap[fmt.Sprintf("AO%s", dataInterval)]    // AO current value
	ta.Value.Oscillators.AO.Prev1 = responseMap[fmt.Sprintf("AO[1]%s", dataInterval)] // AO previous 1 value
	ta.Value.Oscillators.AO.Prev2 = responseMap[fmt.Sprintf("AO[2]%s", dataInterval)] // AO previous 2 value

	// Momentum (10)
	ta.Recommend.Oscillators.Mom = tvMom(responseMap[fmt.Sprintf("Mom%s", dataInterval)], responseMap[fmt.Sprintf("Mom[1]%s", dataInterval)])
	ta.Value.Oscillators.Mom = responseMap[fmt.Sprintf("Mom%s", dataInterval)]

	// MACD Level (12, 26)
	ta.Recommend.Oscillators.MACD = tvMacd(responseMap[fmt.Sprintf("MACD.macd%s", dataInterval)], responseMap[fmt.Sprintf("MACD.signal%s", dataInterval)])
	ta.Value.Oscillators.MACD.Macd = responseMap[fmt.Sprintf("MACD.macd%s", dataInterval)]     // MACD line
	ta.Value.Oscillators.MACD.Signal = responseMap[fmt.Sprintf("MACD.signal%s", dataInterval)] // Signal line

	// Stochastic RSI Fast (3, 3, 14, 14)
	ta.Recommend.Oscillators.StochRSI = tvSimple(responseMap[fmt.Sprintf("Rec.Stoch.RSI%s", dataInterval)])
	ta.Value.Oscillators.StochRSI = responseMap[fmt.Sprintf("Stoch.RSI.K%s", dataInterval)]

	// Williams Percent Range (14)
	ta.Recommend.Oscillators.WR = tvSimple(responseMap[fmt.Sprintf("Rec.WR%s", dataInterval)])
	ta.Value.Oscillators.WR = responseMap[fmt.Sprintf("W.R%s", dataInterval)]

	// Bull Bear Power
	ta.Recommend.Oscillators.BBP = tvSimple(responseMap[fmt.Sprintf("Rec.BBPower%s", dataInterval)])
	ta.Value.Oscillators.BBP = responseMap[fmt.Sprintf("BBPower%s", dataInterval)]

	// Ultimate Oscillator (7, 14, 28)
	ta.Recommend.Oscillators.UO = tvSimple(responseMap[fmt.Sprintf("Rec.UO%s", dataInterval)])
	ta.Value.Oscillators.UO = responseMap[fmt.Sprintf("UO%s", dataInterval)]

	// Moving Averages
	// Exponential Moving Average (EMA)
	ta.Recommend.MovingAverages.EMA10 = tvMa(responseMap[fmt.Sprintf("EMA10%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA10 = responseMap[fmt.Sprintf("EMA10%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA20 = tvMa(responseMap[fmt.Sprintf("EMA20%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA20 = responseMap[fmt.Sprintf("EMA20%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA30 = tvMa(responseMap[fmt.Sprintf("EMA30%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA30 = responseMap[fmt.Sprintf("EMA30%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA50 = tvMa(responseMap[fmt.Sprintf("EMA50%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA50 = responseMap[fmt.Sprintf("EMA50%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA100 = tvMa(responseMap[fmt.Sprintf("EMA100%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA100 = responseMap[fmt.Sprintf("EMA100%s", dataInterval)]

	ta.Recommend.MovingAverages.EMA200 = tvMa(responseMap[fmt.Sprintf("EMA200%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.EMA200 = responseMap[fmt.Sprintf("EMA200%s", dataInterval)]

	// Simple Moving Average (SMA)
	ta.Recommend.MovingAverages.SMA10 = tvMa(responseMap[fmt.Sprintf("SMA10%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA10 = responseMap[fmt.Sprintf("SMA10%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA20 = tvMa(responseMap[fmt.Sprintf("SMA20%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA20 = responseMap[fmt.Sprintf("SMA20%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA30 = tvMa(responseMap[fmt.Sprintf("SMA30%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA30 = responseMap[fmt.Sprintf("SMA30%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA50 = tvMa(responseMap[fmt.Sprintf("SMA50%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA50 = responseMap[fmt.Sprintf("SMA50%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA100 = tvMa(responseMap[fmt.Sprintf("SMA100%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA100 = responseMap[fmt.Sprintf("SMA100%s", dataInterval)]

	ta.Recommend.MovingAverages.SMA200 = tvMa(responseMap[fmt.Sprintf("SMA200%s", dataInterval)], responseMap[fmt.Sprintf("close%s", dataInterval)])
	ta.Value.MovingAverages.SMA200 = responseMap[fmt.Sprintf("SMA200%s", dataInterval)]

	// Ichimoku Base Line (9, 26, 52, 26)
	ta.Recommend.MovingAverages.Ichimoku = tvSimple(responseMap[fmt.Sprintf("Rec.Ichimoku%s", dataInterval)])
	ta.Value.MovingAverages.Ichimoku = responseMap[fmt.Sprintf("Ichimoku.BLine%s", dataInterval)]

	// Volume Weighted Moving Average (20)
	ta.Recommend.MovingAverages.VWMA = tvSimple(responseMap[fmt.Sprintf("Rec.VWMA%s", dataInterval)])
	ta.Value.MovingAverages.VWMA = responseMap[fmt.Sprintf("VWMA%s", dataInterval)]

	// Hull Moving Average (9)
	ta.Recommend.MovingAverages.HullMA = tvSimple(responseMap[fmt.Sprintf("Rec.HullMA9%s", dataInterval)])
	ta.Value.MovingAverages.HullMA = responseMap[fmt.Sprintf("HullMA9%s", dataInterval)]

	// Pivots
	// Pivots - Classic
	ta.Value.Pivots.Classic.Middle = responseMap[fmt.Sprintf("Pivot.M.Classic.Middle%s", dataInterval)]
	ta.Value.Pivots.Classic.R1 = responseMap[fmt.Sprintf("Pivot.M.Classic.R1%s", dataInterval)]
	ta.Value.Pivots.Classic.R2 = responseMap[fmt.Sprintf("Pivot.M.Classic.R2%s", dataInterval)]
	ta.Value.Pivots.Classic.R3 = responseMap[fmt.Sprintf("Pivot.M.Classic.R3%s", dataInterval)]
	ta.Value.Pivots.Classic.S1 = responseMap[fmt.Sprintf("Pivot.M.Classic.S1%s", dataInterval)]
	ta.Value.Pivots.Classic.S2 = responseMap[fmt.Sprintf("Pivot.M.Classic.S2%s", dataInterval)]
	ta.Value.Pivots.Classic.S3 = responseMap[fmt.Sprintf("Pivot.M.Classic.S3%s", dataInterval)]

	// Pivots - Fibonacci
	ta.Value.Pivots.Fibonacci.Middle = responseMap[fmt.Sprintf("Pivot.M.Fibonacci.Middle%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.R1 = responseMap[fmt.Sprintf("Pivot.M.Fibonacci.R1%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.R2 = responseMap[fmt.Sprintf("Pivot.M.Fibonacci.R2%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.R3 = responseMap[fmt.Sprintf("Pivot.M.Fibonacci.R3%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.S1 = responseMap[fmt.Sprintf("Pivot.M.Fibonacci.S1%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.S2 = responseMap[fmt.Sprintf("Pivot.M.Fibonacci.S2%s", dataInterval)]
	ta.Value.Pivots.Fibonacci.S3 = responseMap[fmt.Sprintf("Pivot.M.Fibonacci.S3%s", dataInterval)]

	// Pivots - Camarilla
	ta.Value.Pivots.Camarilla.Middle = responseMap[fmt.Sprintf("Pivot.M.Camarilla.Middle%s", dataInterval)]
	ta.Value.Pivots.Camarilla.R1 = responseMap[fmt.Sprintf("Pivot.M.Camarilla.R1%s", dataInterval)]
	ta.Value.Pivots.Camarilla.R2 = responseMap[fmt.Sprintf("Pivot.M.Camarilla.R2%s", dataInterval)]
	ta.Value.Pivots.Camarilla.R3 = responseMap[fmt.Sprintf("Pivot.M.Camarilla.R3%s", dataInterval)]
	ta.Value.Pivots.Camarilla.S1 = responseMap[fmt.Sprintf("Pivot.M.Camarilla.S1%s", dataInterval)]
	ta.Value.Pivots.Camarilla.S2 = responseMap[fmt.Sprintf("Pivot.M.Camarilla.S2%s", dataInterval)]
	ta.Value.Pivots.Camarilla.S3 = responseMap[fmt.Sprintf("Pivot.M.Camarilla.S3%s", dataInterval)]

	// Pivots - Woodie
	ta.Value.Pivots.Woodie.Middle = responseMap[fmt.Sprintf("Pivot.M.Woodie.Middle%s", dataInterval)]
	ta.Value.Pivots.Woodie.R1 = responseMap[fmt.Sprintf("Pivot.M.Woodie.R1%s", dataInterval)]
	ta.Value.Pivots.Woodie.R2 = responseMap[fmt.Sprintf("Pivot.M.Woodie.R2%s", dataInterval)]
	ta.Value.Pivots.Woodie.R3 = responseMap[fmt.Sprintf("Pivot.M.Woodie.R3%s", dataInterval)]
	ta.Value.Pivots.Woodie.S1 = responseMap[fmt.Sprintf("Pivot.M.Woodie.S1%s", dataInterval)]
	ta.Value.Pivots.Woodie.S2 = responseMap[fmt.Sprintf("Pivot.M.Woodie.S2%s", dataInterval)]
	ta.Value.Pivots.Woodie.S3 = responseMap[fmt.Sprintf("Pivot.M.Woodie.S3%s", dataInterval)]

	// Pivots - Demark
	ta.Value.Pivots.Demark.Middle = responseMap[fmt.Sprintf("Pivot.M.Demark.Middle%s", dataInterval)]
	ta.Value.Pivots.Demark.R1 = responseMap[fmt.Sprintf("Pivot.M.Demark.R1%s", dataInterval)]
	ta.Value.Pivots.Demark.S1 = responseMap[fmt.Sprintf("Pivot.M.Demark.S1%s", dataInterval)]

	// Prices
	ta.Value.Prices.Close = responseMap[fmt.Sprintf("close%s", dataInterval)]
	ta.Value.Prices.High = responseMap[fmt.Sprintf("high%s", dataInterval)]
	ta.Value.Prices.Low = responseMap[fmt.Sprintf("low%s", dataInterval)]

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
