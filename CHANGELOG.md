# Changelog

## v1.4.0

Requires Go 1.27 or later. The module path remains
`github.com/artlevitan/go-tradingview-ta`.

### Added

- `TradingView.GetContext` and `Client.GetContext` for cancellation, deadlines,
  and strict interval validation.
- `RawValues`, `Raw.Lookup`, `Raw.Map`, and `TradingView.HasValue` to distinguish
  unavailable indicator values from real zeroes. Raw JSON retains nulls.
- `SignalUnavailable` (`3`), a separate state outside the rating scale `[-2, 2]`.
- Snapshot metadata: symbol, normalized interval, retrieval time, requested
  candle opening time, and TradingView's update mode.
- `HTTPError` with a status code, bounded response excerpt, and `RetryAfter`.
- Deterministic regression tests, signal boundary tests, cancellation and
  concurrency checks, fuzzing, and opt-in live integration checks.
- Complete public value mapping and recommendation-availability checks,
  rejection tests for every missing scanner field, and mid-response deadlines.
- Live checks for eight venues/providers across all ten supported intervals
  (80 combinations), plus a real nonexistent-symbol rejection.
- CI for Go 1.27 on Linux, macOS, and Windows. Live checks are manually enabled.

### Fixed

- Null moving averages no longer produce false buy recommendations. This
  includes EMA200 and SMA200 on monthly instruments with insufficient history.
- Incomplete and malformed scanner responses return `ErrInvalidResponse`
  instead of becoming zero-valued snapshots. Invalid responses and canceled
  requests leave the previous snapshot unchanged.
- Strict JSON decoding rejects duplicate keys, invalid UTF-8, invalid numeric
  ratings, and candle times outside the RFC 3339 year range.
- Successful responses are limited to 1 MiB; HTTP error excerpts to 4 KiB.
- Scanner URLs preserve existing query parameters and correctly replace
  `symbol` and `fields` without manual URL concatenation.
- Symbol validation rejects empty components, whitespace, and control characters.

### Compatibility

- All previously exported methods, constants, and field types remain available.
  `TradingView` remains comparable and usable as a map key.
- `Get` retains v1.3's documented daily fallback for unknown intervals.
  Use `GetContext` to receive `ErrInvalidInterval` for unsupported intervals.
- Existing numeric `Value` fields remain `float64`; unavailable values remain
  zero there. Check `HasValue` or `Raw.Lookup` before interpreting a number.
- Handle `SignalUnavailable` explicitly; it is neither neutral nor a buy signal.
- `Raw` and `Metadata` are new serialized fields. Prefer comparing `Value` and
  `Recommend` explicitly when comparing analysis results between requests.
- Aggregate ratings still come directly from TradingView. Existing individual
  Stochastic and ADX crossover rules are preserved; documentation and observed
  website behavior differ, so universal formula equivalence is not claimed.

Go module release versions are assigned by Git tags. Publishing v1.4.0 requires
tagging the tested release commit as `v1.4.0` and pushing that tag; changing
`go.mod` changes the required Go version, not the module's release number.
