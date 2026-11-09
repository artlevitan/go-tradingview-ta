// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

// indicatorValues contains only available numbers, keyed without interval suffixes.
type indicatorValues map[string]float64

func (values indicatorValues) recommend(signal int, fields ...string) int {
	for _, field := range fields {
		if _, ok := values[field]; !ok {
			return SignalUnavailable
		}
	}
	return signal
}

func (values indicatorValues) movingAverage(field string) int {
	return values.recommend(movingAverageSignal(values[field], values["close"]), field, "close")
}

func (values indicatorValues) simple(recommendation, field string) int {
	return values.recommend(simpleSignal(values[recommendation]), recommendation, field)
}

func (values indicatorValues) aggregate(field string) int {
	return values.recommend(aggregateSignal(values[field]), field)
}

func (ta *TradingView) populate(values indicatorValues) {
	ta.populateGlobal(values)
	ta.populateOscillators(values)
	ta.populateMovingAverages(values)
	ta.populatePivots(values)
	ta.populatePrices(values)
}

func (ta *TradingView) populateGlobal(values indicatorValues) {
	ta.Recommend.Global.Summary = values.aggregate("Recommend.All")
	ta.Value.Global.Summary = values["Recommend.All"]

	ta.Recommend.Global.Oscillators = values.aggregate("Recommend.Other")
	ta.Value.Global.Oscillators = values["Recommend.Other"]

	ta.Recommend.Global.MA = values.aggregate("Recommend.MA")
	ta.Value.Global.MA = values["Recommend.MA"]
}

func (ta *TradingView) populateOscillators(values indicatorValues) {
	ta.Recommend.Oscillators.RSI = values.recommend(rsiSignal(values["RSI"], values["RSI[1]"]), "RSI", "RSI[1]")
	ta.Value.Oscillators.RSI = values["RSI"]

	ta.Recommend.Oscillators.StochK = values.recommend(stochasticSignal(values["Stoch.K"], values["Stoch.D"], values["Stoch.K[1]"], values["Stoch.D[1]"]), "Stoch.K", "Stoch.D", "Stoch.K[1]", "Stoch.D[1]")
	ta.Value.Oscillators.StochK = values["Stoch.K"]

	ta.Recommend.Oscillators.CCI = values.recommend(cciSignal(values["CCI20"], values["CCI20[1]"]), "CCI20", "CCI20[1]")
	ta.Value.Oscillators.CCI = values["CCI20"]

	ta.Recommend.Oscillators.ADX = values.recommend(adxSignal(values["ADX"], values["ADX+DI"], values["ADX-DI"], values["ADX+DI[1]"], values["ADX-DI[1]"]), "ADX", "ADX+DI", "ADX-DI", "ADX+DI[1]", "ADX-DI[1]")
	ta.Value.Oscillators.ADX.Value = values["ADX"]
	ta.Value.Oscillators.ADX.PlusDI = values["ADX+DI"]
	ta.Value.Oscillators.ADX.MinusDI = values["ADX-DI"]
	ta.Value.Oscillators.ADX.PlusDI1 = values["ADX+DI[1]"]
	ta.Value.Oscillators.ADX.MinusDI1 = values["ADX-DI[1]"]

	ta.Recommend.Oscillators.AO = values.recommend(awesomeOscillatorSignal(values["AO"], values["AO[1]"], values["AO[2]"]), "AO", "AO[1]", "AO[2]")
	ta.Value.Oscillators.AO.Value = values["AO"]
	ta.Value.Oscillators.AO.Prev1 = values["AO[1]"]
	ta.Value.Oscillators.AO.Prev2 = values["AO[2]"]

	ta.Recommend.Oscillators.Mom = values.recommend(momentumSignal(values["Mom"], values["Mom[1]"]), "Mom", "Mom[1]")
	ta.Value.Oscillators.Mom = values["Mom"]

	ta.Recommend.Oscillators.MACD = values.recommend(macdSignal(values["MACD.macd"], values["MACD.signal"]), "MACD.macd", "MACD.signal")
	ta.Value.Oscillators.MACD.Macd = values["MACD.macd"]
	ta.Value.Oscillators.MACD.Signal = values["MACD.signal"]

	ta.Recommend.Oscillators.StochRSI = values.simple("Rec.Stoch.RSI", "Stoch.RSI.K")
	ta.Value.Oscillators.StochRSI = values["Stoch.RSI.K"]

	ta.Recommend.Oscillators.WR = values.simple("Rec.WR", "W.R")
	ta.Value.Oscillators.WR = values["W.R"]

	ta.Recommend.Oscillators.BBP = values.simple("Rec.BBPower", "BBPower")
	ta.Value.Oscillators.BBP = values["BBPower"]

	ta.Recommend.Oscillators.UO = values.simple("Rec.UO", "UO")
	ta.Value.Oscillators.UO = values["UO"]
}

func (ta *TradingView) populateMovingAverages(values indicatorValues) {
	ta.Recommend.MovingAverages.EMA10 = values.movingAverage("EMA10")
	ta.Value.MovingAverages.EMA10 = values["EMA10"]

	ta.Recommend.MovingAverages.EMA20 = values.movingAverage("EMA20")
	ta.Value.MovingAverages.EMA20 = values["EMA20"]

	ta.Recommend.MovingAverages.EMA30 = values.movingAverage("EMA30")
	ta.Value.MovingAverages.EMA30 = values["EMA30"]

	ta.Recommend.MovingAverages.EMA50 = values.movingAverage("EMA50")
	ta.Value.MovingAverages.EMA50 = values["EMA50"]

	ta.Recommend.MovingAverages.EMA100 = values.movingAverage("EMA100")
	ta.Value.MovingAverages.EMA100 = values["EMA100"]

	ta.Recommend.MovingAverages.EMA200 = values.movingAverage("EMA200")
	ta.Value.MovingAverages.EMA200 = values["EMA200"]

	ta.Recommend.MovingAverages.SMA10 = values.movingAverage("SMA10")
	ta.Value.MovingAverages.SMA10 = values["SMA10"]

	ta.Recommend.MovingAverages.SMA20 = values.movingAverage("SMA20")
	ta.Value.MovingAverages.SMA20 = values["SMA20"]

	ta.Recommend.MovingAverages.SMA30 = values.movingAverage("SMA30")
	ta.Value.MovingAverages.SMA30 = values["SMA30"]

	ta.Recommend.MovingAverages.SMA50 = values.movingAverage("SMA50")
	ta.Value.MovingAverages.SMA50 = values["SMA50"]

	ta.Recommend.MovingAverages.SMA100 = values.movingAverage("SMA100")
	ta.Value.MovingAverages.SMA100 = values["SMA100"]

	ta.Recommend.MovingAverages.SMA200 = values.movingAverage("SMA200")
	ta.Value.MovingAverages.SMA200 = values["SMA200"]

	ta.Recommend.MovingAverages.Ichimoku = values.simple("Rec.Ichimoku", "Ichimoku.BLine")
	ta.Value.MovingAverages.Ichimoku = values["Ichimoku.BLine"]

	ta.Recommend.MovingAverages.VWMA = values.simple("Rec.VWMA", "VWMA")
	ta.Value.MovingAverages.VWMA = values["VWMA"]

	ta.Recommend.MovingAverages.HullMA = values.simple("Rec.HullMA9", "HullMA9")
	ta.Value.MovingAverages.HullMA = values["HullMA9"]
}

func (ta *TradingView) populatePivots(values indicatorValues) {
	ta.Value.Pivots.Classic.Middle = values["Pivot.M.Classic.Middle"]
	ta.Value.Pivots.Classic.R1 = values["Pivot.M.Classic.R1"]
	ta.Value.Pivots.Classic.R2 = values["Pivot.M.Classic.R2"]
	ta.Value.Pivots.Classic.R3 = values["Pivot.M.Classic.R3"]
	ta.Value.Pivots.Classic.S1 = values["Pivot.M.Classic.S1"]
	ta.Value.Pivots.Classic.S2 = values["Pivot.M.Classic.S2"]
	ta.Value.Pivots.Classic.S3 = values["Pivot.M.Classic.S3"]

	ta.Value.Pivots.Fibonacci.Middle = values["Pivot.M.Fibonacci.Middle"]
	ta.Value.Pivots.Fibonacci.R1 = values["Pivot.M.Fibonacci.R1"]
	ta.Value.Pivots.Fibonacci.R2 = values["Pivot.M.Fibonacci.R2"]
	ta.Value.Pivots.Fibonacci.R3 = values["Pivot.M.Fibonacci.R3"]
	ta.Value.Pivots.Fibonacci.S1 = values["Pivot.M.Fibonacci.S1"]
	ta.Value.Pivots.Fibonacci.S2 = values["Pivot.M.Fibonacci.S2"]
	ta.Value.Pivots.Fibonacci.S3 = values["Pivot.M.Fibonacci.S3"]

	ta.Value.Pivots.Camarilla.Middle = values["Pivot.M.Camarilla.Middle"]
	ta.Value.Pivots.Camarilla.R1 = values["Pivot.M.Camarilla.R1"]
	ta.Value.Pivots.Camarilla.R2 = values["Pivot.M.Camarilla.R2"]
	ta.Value.Pivots.Camarilla.R3 = values["Pivot.M.Camarilla.R3"]
	ta.Value.Pivots.Camarilla.S1 = values["Pivot.M.Camarilla.S1"]
	ta.Value.Pivots.Camarilla.S2 = values["Pivot.M.Camarilla.S2"]
	ta.Value.Pivots.Camarilla.S3 = values["Pivot.M.Camarilla.S3"]

	ta.Value.Pivots.Woodie.Middle = values["Pivot.M.Woodie.Middle"]
	ta.Value.Pivots.Woodie.R1 = values["Pivot.M.Woodie.R1"]
	ta.Value.Pivots.Woodie.R2 = values["Pivot.M.Woodie.R2"]
	ta.Value.Pivots.Woodie.R3 = values["Pivot.M.Woodie.R3"]
	ta.Value.Pivots.Woodie.S1 = values["Pivot.M.Woodie.S1"]
	ta.Value.Pivots.Woodie.S2 = values["Pivot.M.Woodie.S2"]
	ta.Value.Pivots.Woodie.S3 = values["Pivot.M.Woodie.S3"]

	ta.Value.Pivots.Demark.Middle = values["Pivot.M.Demark.Middle"]
	ta.Value.Pivots.Demark.R1 = values["Pivot.M.Demark.R1"]
	ta.Value.Pivots.Demark.S1 = values["Pivot.M.Demark.S1"]
}

func (ta *TradingView) populatePrices(values indicatorValues) {
	ta.Value.Prices.Close = values["close"]
	ta.Value.Prices.High = values["high"]
	ta.Value.Prices.Low = values["low"]
}

// aggregateSignal converts TradingView's aggregate score into a public signal.
func aggregateSignal(v float64) int {
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

// rsiSignal converts RSI values into a normalized recommendation signal.
func rsiSignal(rsi, rsi1 float64) int {
	switch {
	case rsi < 30 && rsi1 < rsi:
		return SignalBuy
	case rsi > 70 && rsi1 > rsi:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// stochasticSignal converts stochastic values into a normalized recommendation signal.
func stochasticSignal(k, d, k1, d1 float64) int {
	switch {
	case k < 20 && d < 20 && k > d && k1 < d1:
		return SignalBuy
	case k > 80 && d > 80 && k < d && k1 > d1:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// cciSignal converts CCI values into a normalized recommendation signal.
func cciSignal(cci20, cci201 float64) int {
	switch {
	case cci20 < -100 && cci20 > cci201:
		return SignalBuy
	case cci20 > 100 && cci20 < cci201:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// adxSignal converts ADX values into a normalized recommendation signal.
func adxSignal(adx, adxpdi, adxndi, adxpdi1, adxndi1 float64) int {
	switch {
	case adx > 20 && adxpdi1 < adxndi1 && adxpdi > adxndi:
		return SignalBuy
	case adx > 20 && adxpdi1 > adxndi1 && adxpdi < adxndi:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// awesomeOscillatorSignal converts Awesome Oscillator values into a normalized recommendation signal.
func awesomeOscillatorSignal(ao, ao1, ao2 float64) int {
	switch {
	case (ao > 0 && ao1 < 0) || (ao > 0 && ao1 > 0 && ao > ao1 && ao2 > ao1):
		return SignalBuy
	case (ao < 0 && ao1 > 0) || (ao < 0 && ao1 < 0 && ao < ao1 && ao2 < ao1):
		return SignalSell
	default:
		return SignalNeutral
	}
}

// momentumSignal converts momentum values into a normalized recommendation signal.
func momentumSignal(mom, mom1 float64) int {
	switch {
	case mom > mom1:
		return SignalBuy
	case mom < mom1:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// macdSignal converts MACD values into a normalized recommendation signal.
func macdSignal(macd, s float64) int {
	switch {
	case macd > s:
		return SignalBuy
	case macd < s:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// simpleSignal converts a TradingView simple recommendation value into a signal.
func simpleSignal(v float64) int {
	switch {
	case v == 1:
		return SignalBuy
	case v == -1:
		return SignalSell
	default:
		return SignalNeutral
	}
}

// movingAverageSignal converts a moving average and close price into a signal.
func movingAverageSignal(ma, close float64) int {
	switch {
	case ma < close:
		return SignalBuy
	case ma > close:
		return SignalSell
	default:
		return SignalNeutral
	}
}
