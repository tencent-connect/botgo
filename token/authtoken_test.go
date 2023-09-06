package token

import "testing"

func Test_getTokenTTL(t *testing.T) {
	type args struct {
		tokenTTL int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "getTokenTTL-7200",
			args: args{tokenTTL: 7200},
			want: 7170,
		},
		{
			name: "getTokenTTL-60",
			args: args{tokenTTL: 60},
			want: 30,
		},
		{
			name: "getTokenTTL-30",
			args: args{tokenTTL: 30},
			want: minTimer,
		},
		{
			name: "getTokenTTL-10",
			args: args{tokenTTL: 10},
			want: minTimer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTokenTTL(tt.args.tokenTTL); got > tt.want {
				t.Errorf("getTokenTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}
