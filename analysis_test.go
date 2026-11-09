// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tradingview

import (
	"math"
	"testing"
)

func TestAggregateSignalBoundaries(t *testing.T) {
	for _, tt := range []struct {
		value float64
		want  int
	}{
		{-1, SignalStrongSell},
		{math.Nextafter(-0.5, -1), SignalStrongSell},
		{-0.5, SignalSell},
		{math.Nextafter(-0.1, -1), SignalSell},
		{-0.1, SignalNeutral},
		{0, SignalNeutral},
		{0.1, SignalNeutral},
		{math.Nextafter(0.1, 1), SignalBuy},
		{0.5, SignalBuy},
		{math.Nextafter(0.5, 1), SignalStrongBuy},
		{1, SignalStrongBuy},
	} {
		if got := aggregateSignal(tt.value); got != tt.want {
			t.Errorf("aggregateSignal(%g) = %d, want %d", tt.value, got, tt.want)
		}
	}
}

func TestIndicatorSignals(t *testing.T) {
	for _, tt := range []struct {
		name string
		got  int
		want int
	}{
		{"RSI oversold rising", rsiSignal(29, 28), SignalBuy},
		{"RSI overbought falling", rsiSignal(71, 72), SignalSell},
		{"RSI lower threshold", rsiSignal(30, 29), SignalNeutral},
		{"RSI upper threshold", rsiSignal(70, 71), SignalNeutral},
		{"RSI oversold falling", rsiSignal(28, 29), SignalNeutral},
		{"RSI overbought rising", rsiSignal(72, 71), SignalNeutral},
		{"RSI unchanged", rsiSignal(29, 29), SignalNeutral},
		{"stochastic bullish crossover", stochasticSignal(19, 18, 17, 18), SignalBuy},
		{"stochastic bearish crossover", stochasticSignal(81, 82, 83, 82), SignalSell},
		{"stochastic no bullish crossover", stochasticSignal(19, 18, 18, 17), SignalNeutral},
		{"stochastic no bearish crossover", stochasticSignal(81, 82, 82, 83), SignalNeutral},
		{"stochastic lower threshold", stochasticSignal(20, 18, 17, 18), SignalNeutral},
		{"stochastic upper threshold", stochasticSignal(80, 82, 83, 82), SignalNeutral},
		{"CCI oversold rising", cciSignal(-101, -102), SignalBuy},
		{"CCI overbought falling", cciSignal(101, 102), SignalSell},
		{"CCI lower threshold", cciSignal(-100, -101), SignalNeutral},
		{"CCI upper threshold", cciSignal(100, 101), SignalNeutral},
		{"CCI unchanged", cciSignal(101, 101), SignalNeutral},
		{"ADX bullish DI crossover", adxSignal(21, 30, 20, 19, 21), SignalBuy},
		{"ADX bearish DI crossover", adxSignal(21, 20, 30, 31, 29), SignalSell},
		{"ADX threshold", adxSignal(20, 30, 20, 19, 21), SignalNeutral},
		{"ADX no crossover", adxSignal(30, 30, 20, 29, 19), SignalNeutral},
		{"ADX equal DI", adxSignal(30, 20, 20, 19, 21), SignalNeutral},
		{"AO crosses positive", awesomeOscillatorSignal(1, -1, -2), SignalBuy},
		{"AO crosses negative", awesomeOscillatorSignal(-1, 1, 2), SignalSell},
		{"AO positive saucer", awesomeOscillatorSignal(3, 1, 2), SignalBuy},
		{"AO negative saucer", awesomeOscillatorSignal(-3, -1, -2), SignalSell},
		{"AO positive falling", awesomeOscillatorSignal(1, 2, 3), SignalNeutral},
		{"AO negative rising", awesomeOscillatorSignal(-1, -2, -3), SignalNeutral},
		{"AO zero", awesomeOscillatorSignal(0, 0, 0), SignalNeutral},
		{"momentum rises", momentumSignal(-1, -2), SignalBuy},
		{"momentum falls", momentumSignal(1, 2), SignalSell},
		{"momentum unchanged", momentumSignal(0, 0), SignalNeutral},
		{"MACD above signal", macdSignal(-1, -2), SignalBuy},
		{"MACD below signal", macdSignal(1, 2), SignalSell},
		{"MACD equal signal", macdSignal(0, 0), SignalNeutral},
		{"simple buy", simpleSignal(1), SignalBuy},
		{"simple sell", simpleSignal(-1), SignalSell},
		{"simple neutral", simpleSignal(0), SignalNeutral},
		{"MA below price", movingAverageSignal(10, 11), SignalBuy},
		{"MA above price", movingAverageSignal(11, 10), SignalSell},
		{"MA equals price", movingAverageSignal(10, 10), SignalNeutral},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("signal = %d, want %d", tt.got, tt.want)
			}
		})
	}
}
