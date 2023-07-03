package util

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "ShouldError_InvalidInput",
			args:    args{str: "invalid"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "ShouldError_InvalidHour",
			args:    args{str: "xx:xx:xx.xxx"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "ShouldError_InvalidMinute",
			args:    args{str: "10:xx:xx.xxx"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "ShouldError_InvalidSecond",
			args:    args{str: "10:20:xx.xxx"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "ShouldError_InvalidMilliSecond",
			args:    args{str: "10:20:00.xxx"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "ShouldError_InvalidMilliSecondValue",
			args:    args{str: "10:20:00.78192487912847124"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Valid",
			args:    args{str: "00:01:00.0"},
			want:    time.Minute,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
