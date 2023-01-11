package tradingview

import (
	"reflect"
	"testing"
)

func TestTradingViewData(t *testing.T) {
	type args struct {
		symbol   string
		interval string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int
		wantErr bool
	}{
		{name: "BINANCE:BTCUSDT", args: args{symbol: "BINANCE:BTCUSDT", interval: interval15min}, want: map[string]int{}, wantErr: false},
		{name: "BINANCE:ETCUSDT", args: args{symbol: "BINANCE:ETCUSDT", interval: interval5min}, want: map[string]int{}, wantErr: false},
		{name: "BINANCE:ETCUSDT", args: args{symbol: "BINANCE:ETCUSDT", interval: interval4hour}, want: map[string]int{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TradingViewData(tt.args.symbol, tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("TradingViewData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TradingViewData() = %v, want %v", got, tt.want)
			}
		})
	}
}
