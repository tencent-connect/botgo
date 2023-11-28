package token

import (
	"testing"
)

func Test_getTokenTTL(t *testing.T) {
	type args struct {
		tokenTTL float64
	}
	tests := []struct {
		name string
		args args
		want float64
		max  float64
		min  float64
	}{
		{
			name: "less than preserveTokenTTL",
			args: args{tokenTTL: preserveTokenTTL - 1},
			want: preserveTokenTTL - 1,
		},
		{
			name: "between preserveTokenTTL and preserveTokenTTL+minTimeGap",
			args: args{tokenTTL: preserveTokenTTL + minTimeGap - 1},
			want: minTimeGap,
		},
		{
			name: "greater than preserveTokenTTL+minTimeGap",
			args: args{tokenTTL: preserveTokenTTL + minTimeGap + 1},
			want: minTimeGap + 1,
		},
		{
			name: "greater than preserveTokenTTL+minTimeGap",
			args: args{tokenTTL: preserveTokenTTL + minTimeGap + 1},
			want: minTimeGap + 1,
		},
		{
			name: "greater than preserveTokenTTL+randUpperLimit",
			args: args{tokenTTL: preserveTokenTTL + randTimeUpperLimit + 1},
			max:  randTimeUpperLimit + 1,
			min:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTokenTTL(tt.args.tokenTTL); got != tt.want && !(got >= tt.min && got <= tt.max) {
				t.Errorf("getTokenTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}
