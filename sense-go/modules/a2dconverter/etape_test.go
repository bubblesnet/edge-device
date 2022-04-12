package a2dconverter

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_etapeInchesFromVolts(t *testing.T) {

	type args struct {
		voltage    float64
		slope      float64
		yintercept float64
	}
	tests := []struct {
		name       string
		args       args
		wantInches float64
		wantError  bool
	}{
		// TODO: Add test cases.
		{name: "Bottom", args: args{voltage: MinVoltage, slope: Etape_slope, yintercept: Etape_y_intercept}, wantInches: 1.49, wantError: false},
		{name: "Low", args: args{voltage: 1.77, slope: Etape_slope, yintercept: Etape_y_intercept}, wantInches: 2.85, wantError: false},
		{name: "Medium", args: args{voltage: 2.09, slope: Etape_slope, yintercept: Etape_y_intercept}, wantInches: 6.49, wantError: false},
		{name: "High", args: args{voltage: 2.49, slope: Etape_slope, yintercept: Etape_y_intercept}, wantInches: 11.05, wantError: false},
		{name: "Top", args: args{voltage: MaxVoltage, slope: Etape_slope, yintercept: Etape_y_intercept}, wantInches: 12.87, wantError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := etapeInchesFromVolts(tt.args.voltage, tt.args.slope, tt.args.yintercept)
			s := fmt.Sprintf("%.2f", in)
			gotInches, _ := strconv.ParseFloat(s, 64)
			if gotInches != tt.wantInches {
				t.Errorf("etapeInchesFromVolts() = %v, want %v tt.wantError %v", gotInches, tt.wantInches, tt.wantError)
			}
		})
	}
}
