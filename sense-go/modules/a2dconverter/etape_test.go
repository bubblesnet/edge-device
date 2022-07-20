/*
 * Copyright (c) John Rodley 2022.
 * SPDX-FileCopyrightText:  John Rodley 2022.
 * SPDX-License-Identifier: MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the
 * Software without restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the
 * following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
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
