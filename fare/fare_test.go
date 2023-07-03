package fare

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func Test_fare_calculateFare(t *testing.T) {
	type fields struct {
		currentElapsedTime time.Duration
		currentFare        int64
		currentDistance    float64
		records            []fareRecord
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
				records:            tt.fields.records,
			}
			f.calculateFare()
			require.Equal(t, tt.expected, f.currentFare)
		})
	}
}

func Test_fare_ShowRecords(t *testing.T) {
	type fields struct {
		currentElapsedTime time.Duration
		currentFare        int64
		currentDistance    float64
		records            []fareRecord
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Success_ShowRecords",
			fields: fields{
				currentElapsedTime: time.Minute,
				currentFare:        400,
				currentDistance:    300,
				records: []fareRecord{
					{
						Timestamp:  "00:00:00.000 ",
						Mileage:    0,
						MileageStr: "0.0",
						Diff:       0,
					},
					{
						Timestamp:  "00:00:01.040 ",
						Mileage:    300,
						MileageStr: "300.0",
						Diff:       300,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fare{
				currentElapsedTime: tt.fields.currentElapsedTime,
				currentFare:        tt.fields.currentFare,
				currentDistance:    tt.fields.currentDistance,
				records:            tt.fields.records,
			}
			f.ShowRecords()
		})
	}
}

func Test_fare_parseInput(t *testing.T) {
	type fields struct {
		currentElapsedTime time.Duration
		currentFare        int64
		currentDistance    float64
		records            []fareRecord
	}
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "ShouldError_ImproperFormat_Invalid",
			fields:  fields{},
			args:    args{input: "asd"},
			wantErr: true,
		},
		{
			name:    "ShouldError_ImproperFormat_InvalidDuration",
			fields:  fields{},
			args:    args{input: "asd asd"},
			wantErr: true,
		},
		{
			name: "ShouldError_ReceivePasTime",
			fields: fields{
				currentElapsedTime: 5 * time.Minute,
			},
			args:    args{input: "00:01:00.123 480.9"},
			wantErr: true,
		},
		{
			name: "ShouldError_TooLongInterval",
			fields: fields{
				currentElapsedTime: time.Minute,
			},
			args:    args{input: "00:10:00.123 480.9"},
			wantErr: true,
		},
		{
			name:    "ShouldError_ImproperFormat_InvalidDistance",
			fields:  fields{},
			args:    args{input: "00:00:00.00 asd"},
			wantErr: true,
		},
		{
			name:    "ShouldError_InvalidStart",
			fields:  fields{},
			args:    args{input: "00:00:00.0 1.0"},
			wantErr: true,
		},
		{
			name: "Success",
			fields: fields{
				currentElapsedTime: time.Minute,
				records: []fareRecord{
					{
						Timestamp: "00:00:00.0",
					},
				},
			},
			args:    args{input: "00:01:00.123 480.9"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fare{
				currentElapsedTime: tt.fields.currentElapsedTime,
				currentFare:        tt.fields.currentFare,
				currentDistance:    tt.fields.currentDistance,
				records:            tt.fields.records,
			}
			if err := f.parseInput(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_fare_Calculate(t *testing.T) {
	type fields struct {
		currentElapsedTime time.Duration
		currentFare        int64
		currentDistance    float64
		records            []fareRecord
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
			name:    "ShouldError_EndButNotEnoughData",
			fields:  fields{},
			args:    args{input: "end"},
			want:    0,
			wantErr: true,
		},
		{
			name: "ShouldError_EndButZeroMileage",
			fields: fields{
				records: []fareRecord{
					{},
					{},
				},
			},
			args:    args{input: "end"},
			want:    0,
			wantErr: true,
		},
		{
			name: "Success_EndResult",
			fields: fields{
				records: []fareRecord{
					{},
					{},
				},
				currentDistance: 10.0,
				currentFare:     10,
			},
			args:    args{input: "end"},
			want:    10,
			wantErr: false,
		},
		{
			name:    "ShouldError_FailedParseInput",
			fields:  fields{},
			args:    args{input: "00:01:00.123 480.x"},
			want:    0,
			wantErr: true,
		},
		{
			name: "Success",
			fields: fields{
				records: []fareRecord{
					{},
				},
			},
			args:    args{input: "00:01:00.123 480.9"},
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
				records:            tt.fields.records,
			}
			got, err := f.Calculate(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calculate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTaxiFare(t *testing.T) {
	tests := []struct {
		name string
		want TaxiFare
	}{
		{
			name: "Success",
			want: NewTaxiFare(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaxiFare(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaxiFare() = %v, want %v", got, tt.want)
			}
		})
	}
}
