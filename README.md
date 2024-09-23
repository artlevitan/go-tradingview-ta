# go-tradingview-ta [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An unofficial Go API simple wrapper to retrieve technical analysis from TradingView.

<img src="/editor/images/promo.jpg" alt="TradingView Go" style="max-width:100%">

## Installation

```bash
go get github.com/artlevitan/go-tradingview-ta
```

### Predefined constants

```go
const (
	// Intervals
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
)
```

## Example

```go
package main

import (
	"fmt"

	tradingview "github.com/artlevitan/go-tradingview-ta"
)

const SYMBOL = "BINANCE:BTCUSDT" // https://www.tradingview.com/symbols/BTCUSDT/technicals/

func main() {
	var ta tradingview.TradingView
	
	// Fetch data for the specified symbol at a 4-hour interval
	err := ta.Get(SYMBOL, tradingview.Interval4Hour)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// Get the summary trading recommendation
	recSummary := ta.Recommend.Global.Summary

	// Print the recommendation based on the signal
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

	// Print the latest closing price
	clPrice := ta.Value.Prices.Close
	fmt.Println("Closing price:", clPrice)
}
```

## License

MIT
