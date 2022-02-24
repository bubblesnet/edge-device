package a2dconverter

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_etapeInchesFromVolts(t *testing.T) {
	type args struct {
		Vdd           float64
		voltage       float64
		minResistance float64
		maxResistance float64
	}
	tests := []struct {
		name       string
		args       args
		wantInches float64
		wantError  bool
	}{
		// TODO: Add test cases.
		{name: "OutOfRangeLow", args: args{3.325, -1.0, 400.0, 2200.0}, wantInches: 1.0, wantError: true},
		{name: "OutOfRangeHigh", args: args{3.325, 10.0, 400.0, 2200.0}, wantInches: 1.0, wantError: true},
		{name: "Bottom", args: args{3.325, 3.325 / 2.0, 400.0, 2200.0}, wantInches: 1.28, wantError: false},
		{name: "Low", args: args{3.325, 1.9, 400.0, 2200.0}, wantInches: 2.92, wantError: false},
		{name: "Medium", args: args{3.325, 2.4, 400.0, 2200.0}, wantInches: 6.38, wantError: false},
		{name: "High", args: args{3.325, 3.0, 400.0, 2200.0}, wantInches: 10.53, wantError: false},
		{name: "Top", args: args{3.325, 3.325, 400.0, 2200.0}, wantInches: 12.78, wantError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in, err := etapeInchesFromVolts(tt.args.Vdd, tt.args.voltage, tt.args.minResistance, tt.args.maxResistance, 1.0, 12.78)
			if (err != nil) != tt.wantError {
				t.Errorf("error %v wantError %v", err, tt.wantError)
			} else {
				s := fmt.Sprintf("%.2f", in)
				gotInches, _ := strconv.ParseFloat(s, 64)
				if err == nil && gotInches != tt.wantInches {
					t.Errorf("etapeInchesFromVolts() = %v, want %v err %v tt.wantError %v", gotInches, tt.wantInches, err, tt.wantError)
				}
			}
		})
	}
}
func Test_etapeOhmsToInches(t *testing.T) {
	type args struct {
		ohms float64
	}
	tests := []struct {
		name       string
		args       args
		wantInches float64
	}{
		// TODO: Add test cases.
		{name: "Top", args: args{ohms: 2200.0}, wantInches: 1.28},
		{name: "High", args: args{ohms: 1400.0}, wantInches: 6.39},
		{name: "Middle", args: args{ohms: 1100.0}, wantInches: 8.31},
		{name: "Low", args: args{ohms: 900.0}, wantInches: 9.58},
		{name: "Bottom", args: args{ohms: 400.0}, wantInches: 12.78},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := fmt.Sprintf("%.2f", etapeOhmsToInches(tt.args.ohms))
			gotInches, _ := strconv.ParseFloat(s, 64)
			if gotInches != tt.wantInches {
				t.Errorf("etapeOhmsToInches() = %v, want %v", gotInches, tt.wantInches)
			}
		})
	}
}

func Test_etapeVoltageToOhms(t *testing.T) {
	type args struct {
		Vdd           float64
		voltage       float64
		minResistance float64
		maxResistance float64
	}
	tests := []struct {
		name     string
		args     args
		wantOhms float64
	}{
		// TODO: Add test cases.
		{name: "Top", args: args{3.325, 3.325, 400.0, 2200.0}, wantOhms: 400.0},
		{name: "High", args: args{3.325, 3.0, 400.0, 2200.0}, wantOhms: 751.88},
		{name: "Middle", args: args{3.325, 2.5, 400.0, 2200.0}, wantOhms: 1293.23},
		{name: "Low", args: args{3.325, 2.0, 400.0, 2200.0}, wantOhms: 1834.59},
		{name: "Bottom", args: args{3.325, 3.325 / 2.0, 400.0, 2200.0}, wantOhms: 2200.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ohms, _ := etapeVoltageToOhms(tt.args.Vdd, tt.args.voltage, tt.args.minResistance, tt.args.maxResistance)
			s := fmt.Sprintf("%.2f", ohms)
			gotOhms, _ := strconv.ParseFloat(s, 64)
			if gotOhms != tt.wantOhms {
				t.Errorf("%v etapeVoltageToOhms() = %v, want %v", tt.args, gotOhms, tt.wantOhms)
			}
		})
	}
}
