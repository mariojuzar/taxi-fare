package fare

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_fare_CalculateFareMeter(t *testing.T) {
	type fields struct {
		currentElapsedTime time.Duration
		currentFare        int64
		currentDistance    float64
		countHistory       uint
	}
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "ShouldError_BlankInput",
			fields:  fields{},
			args:    args{input: ""},
			want:    0,
			wantErr: true,
		},
		{
			name:    "ShouldError_Finished_NotEnoughData",
			fields:  fields{},
			args:    args{input: "finished"},
			want:    0,
			wantErr: true,
		},
		{
			name: "ShouldError_Finished_ZeroDistance",
			fields: fields{
				countHistory: 2,
			},
			args:    args{input: "finished"},
			want:    0,
			wantErr: true,
		},
		{
			name: "Valid_Finished",
			fields: fields{
				countHistory:    2,
				currentDistance: 1,
				currentFare:     10,
			},
			args:    args{input: "finished"},
			want:    10,
			wantErr: false,
		},
		{
			name:    "ShouldError_InvalidInput",
			fields:  fields{},
			args:    args{input: "invalid"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "ShouldError_FailedParseDuration",
			fields:  fields{},
			args:    args{input: "invalid invalid"},
			want:    0,
			wantErr: true,
		},
		{
			name: "ShouldError_PastTime",
			fields: fields{
				currentElapsedTime: time.Minute,
			},
			args:    args{input: "00:00:01 400"},
			want:    0,
			wantErr: true,
		},
		{
			name: "ShouldError_IntervalTooLong",
			fields: fields{
				currentElapsedTime: time.Minute,
			},
			args:    args{input: "00:10:01 400"},
			want:    0,
			wantErr: true,
		},
		{
			name: "ShouldError_FailedParseDistance",
			fields: fields{
				currentElapsedTime: time.Minute,
			},
			args:    args{input: "00:02:01 xxx"},
			want:    0,
			wantErr: true,
		},
		{
			name: "Success",
			fields: fields{
				currentElapsedTime: time.Minute,
			},
			args:    args{input: "00:02:01 480.9"},
			want:    400,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fare{
				currentElapsedTime: tt.fields.currentElapsedTime,
				currentFare:        tt.fields.currentFare,
				currentDistance:    tt.fields.currentDistance,
				countHistory:       tt.fields.countHistory,
			}
			got, err := f.CalculateFareMeter(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateFareMeter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateFareMeter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fare_calculateFare(t *testing.T) {
	type fields struct {
		currentElapsedTime time.Duration
		currentFare        int64
		currentDistance    float64
		countHistory       uint
	}
	tests := []struct {
		name     string
		fields   fields
		expected int64
	}{
		{
			name:     "CalculateFareWithRule",
			fields:   fields{currentDistance: 10350},
			expected: 1360,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fare{
				currentElapsedTime: tt.fields.currentElapsedTime,
				currentFare:        tt.fields.currentFare,
				currentDistance:    tt.fields.currentDistance,
				countHistory:       tt.fields.countHistory,
			}
			f.calculateFare()
			require.Equal(t, tt.expected, f.currentFare)
		})
	}
}
