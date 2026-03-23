// Copyright 2022-2026 The go-tradingview-ta Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
// Package tradingview provides an unofficial Go client for TradingView's public
// scanner endpoint.
//
// The package exposes normalized recommendation signals together with the raw
// numeric values returned by TradingView for a symbol and interval.
//
// Most callers can use (*TradingView).Get, which delegates to DefaultClient.
// Callers that need a custom HTTP client or endpoint can use Client directly.
//
// The zero values of TradingView and Client are ready to use.
package tradingview
