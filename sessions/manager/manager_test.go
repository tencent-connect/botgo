package manager

import (
	"testing"
	"time"
)

func Test_calcInterval(t *testing.T) {
	type args struct {
		maxConcurrency uint32
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{"c1", args{maxConcurrency: 1}, concurrencyTimeWindowSec * time.Second},
		{"c3", args{maxConcurrency: 3}, 1 * time.Second},
		{"c5", args{maxConcurrency: 5}, 1 * time.Second},
		{"c10", args{maxConcurrency: 10}, 1 * time.Second},
		{"c100", args{maxConcurrency: 100}, 1 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcInterval(tt.args.maxConcurrency); got != tt.want {
				t.Errorf("CalcInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
