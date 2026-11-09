// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
// Package tradingview provides an unofficial Go client for TradingView's public
// scanner endpoint.
//
// The package exposes normalized recommendation signals together with the raw
// numeric values returned by TradingView for a symbol and interval.
// Raw preserves null indicator values; HasValue distinguishes them from real
// zeroes. Recommendations with unavailable inputs are SignalUnavailable.
// Metadata exposes the requested candle time and the scanner's update mode.
//
// Most callers should use [TradingView.GetContext], which delegates to
// [DefaultClient] and supports cancellation, deadlines, and interval validation.
// [TradingView.Get] preserves the v1.3 fallback from unknown intervals to daily
// data. Callers that need a custom HTTP client or endpoint can use [Client].
//
// The zero values of TradingView and Client are ready to use.
package tradingview
