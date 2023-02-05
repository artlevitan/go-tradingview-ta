package tradingview

import (
	"testing"
)

func TestTradingView_Get(t *testing.T) {
	type fields struct {
		Recommend struct {
			Summary     int
			Oscillators int
			MA          int
		}
		Oscillators struct {
			RSI      int
			StochK   int
			CCI      int
			ADX      int
			AO       int
			Mom      int
			MACD     int
			StochRSI int
			WR       int
			BBP      int
			UO       int
		}
		MovingAverages struct {
			EMA10    int
			SMA10    int
			EMA20    int
			SMA20    int
			EMA30    int
			SMA30    int
			EMA50    int
			SMA50    int
			EMA100   int
			SMA100   int
			EMA200   int
			SMA200   int
			Ichimoku int
			VWMA     int
			HullMA   int
		}
	}
	type args struct {
		symbol   string
		interval string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "BINANCE:BTCUSDT - 15min", fields: fields{}, args: args{symbol: "BINANCE:BTCUSDT", interval: Interval15min}, wantErr: false},
		{name: "BINANCE:BTCUSDT - 1hour", fields: fields{}, args: args{symbol: "BINANCE:BTCUSDT", interval: Interval1hour}, wantErr: false},
		{name: "BINANCE:BTCUSDT - 4hour", fields: fields{}, args: args{symbol: "BINANCE:BTCUSDT", interval: Interval4hour}, wantErr: false},
		{name: "BINANCE:BTCUSDT - 1day", fields: fields{}, args: args{symbol: "BINANCE:BTCUSDT", interval: Interval1day}, wantErr: false},
		{name: "BINANCE:BTCUSDT - 1week", fields: fields{}, args: args{symbol: "BINANCE:BTCUSDT", interval: Interval1week}, wantErr: false},
		{name: "BINANCE:BTCUSDT - 1month", fields: fields{}, args: args{symbol: "BINANCE:BTCUSDT", interval: Interval1month}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TradingView{
				Recommend:      tt.fields.Recommend,
				Oscillators:    tt.fields.Oscillators,
				MovingAverages: tt.fields.MovingAverages,
			}
			if err := tr.Get(tt.args.symbol, tt.args.interval); (err != nil) != tt.wantErr {
				t.Errorf("TradingView.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
