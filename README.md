# go-tradingview-ta [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
An unofficial Go API simple wrapper to retrieve technical analysis from TradingView.

## Example
```go
package main

import (
	"fmt"

	tradingview "github.com/artlevitan/go-tradingview-ta"
)

func main() {
	var tv tradingview.TVData
	err := tv.Get("BINANCE:BTCUSDT", tradingview.Interval4hour)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", tv) // Full Data

	// Get the value by key
	recSummary := tv.Recommend.Summary // Summary
	fmt.Println(recSummary)

	recOsc := tv.Recommend.Oscillators // Oscillators
	fmt.Println(recOsc)

	recMA := tv.Recommend.MA // Moving Averages
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