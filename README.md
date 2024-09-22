# go-tradingview-ta [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An unofficial Go API simple wrapper to retrieve technical analysis from TradingView.

<img src="/editor/images/promo.jpg" alt="Go TradingView" style="max-width:100%">

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

func main() {
	var ta tradingview.TradingView
	err := ta.Get("BINANCE:BTCUSDT", tradingview.Interval4Hour)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", ta) // Full Data

	// Summary recommendation
	recSummary := ta.Recommend.Global.Summary
	fmt.Println(recSummary)

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

	// Ichimoku Base Line (9, 26, 52, 26)
	ichimoku := ta.Recommend.MovingAverages.Ichimoku 
	fmt.Println(ichimoku)
}

```

## License

MIT
