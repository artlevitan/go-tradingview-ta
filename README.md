# go-tradingview-ta [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An unofficial Go API simple wrapper to retrieve technical analysis from TradingView.

<img src="/editor/images/promo.jpg" alt="Go TradingView" style="max-width:100%">

### Predefined constants

```go
const (
    // Intervals
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

MIT
