package main

import "testing"

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name           string
		temperatureStr string
		want           float64
		wantErr        bool
	}{
		{"Positive no rounding", "23.50", 23.5, false},
		{"Positive with rounding", "10.76", 10.8, false},
		{"Negative no rounding", "-5.28", -5.2, false},
		{"Negative with rounding", "-2.01", -2.0, false},
		{"Just under negative", "0.99", 1.0, false},
		{"Invalid format", "invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloat(tt.temperatureStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
