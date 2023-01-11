# go-tradingview-ta [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
An unofficial Go API simple wrapper to retrieve technical analysis from TradingView.

### Predefined constants
```go
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
	SignalStrongBuy  = 2  // STRONG_BUY
	SignalBuy        = 1  // BUY
	SignalNeutral    = 0  // NEUTRAL
	SignalSell       = -1 // SELL
	SignalStrongSell = -2 // STRONG_SELL
)

```

## Example
```go
package main

import (
	"fmt"

	tradingview "github.com/artlevitan/go-tradingview-ta"
)

func main() {
	var ta tradingview.TradingView
	err := ta.Get("BINANCE:BTCUSDT", tradingview.Interval4hour)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", ta) // Full Data

	// Get the value by key
	recSummary := ta.Recommend.Summary // Summary
	fmt.Println(recSummary)

	recOsc := ta.Recommend.Oscillators // Oscillators
	fmt.Println(recOsc)

	recMA := ta.Recommend.MA // Moving Averages
	fmt.Println(recMA)

	// Text recommendation
	switch recSummary {
	case tradingview.SignalStrongSell:
		fmt.Println("STRONG_SELL")
	case tradingview.SignalSell:
		fmt.Println("SELL")
	case tradingview.SignalNeutral:
		fmt.Println("NEUTRAL")
	case tradingview.SignalBuy:
		fmt.Println("BUY")
	case tradingview.SignalStrongBuy:
		fmt.Println("STRONG_BUY")
	default:
		fmt.Println("An error has occurred")
	}
}

```

## License
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.